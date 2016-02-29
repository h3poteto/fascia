package list

import (
	"../../config"
	"../../modules/hub"
	"../../modules/logging"
	"../db"
	"../list_option"
	"../repository"
	"../task"
	"database/sql"
)

type List interface {
	Save() bool
}

type ListStruct struct {
	ID           int64
	ProjectID    int64
	UserID       int64
	Title        sql.NullString
	ListTasks    []*task.TaskStruct
	Color        sql.NullString
	ListOptionID sql.NullInt64
	database     db.DB
}

func NewList(id int64, projectID int64, userID int64, title string, color string, optionID sql.NullInt64) *ListStruct {
	if projectID == 0 {
		return nil
	}
	nullTitle := sql.NullString{String: title, Valid: true}
	nullColor := sql.NullString{String: color, Valid: true}

	list := &ListStruct{ID: id, ProjectID: projectID, UserID: userID, Title: nullTitle, Color: nullColor, ListOptionID: optionID}
	list.Initialize()
	return list
}

func FindList(projectID int64, listID int64) *ListStruct {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	var id, userID int64
	var title, color sql.NullString
	var optionID sql.NullInt64
	rows, _ := table.Query("select id, project_id, user_id, title, color, list_option_id from lists where id = ? AND project_id = ?;", listID, projectID)
	for rows.Next() {
		err := rows.Scan(&id, &projectID, &userID, &title, &color, &optionID)
		if err != nil {
			logging.SharedInstance().MethodInfo("List", "FindList", true).Panic(err)
		}
	}
	if id != listID {
		logging.SharedInstance().MethodInfo("list", "FindList", true).Errorf("cannot find list or project did not contain list: %v", listID)
		return nil
	} else {
		list := NewList(id, projectID, userID, title.String, color.String, optionID)
		return list
	}

}

func (u *ListStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}

func (u *ListStruct) Save(repo *repository.RepositoryStruct, OauthToken *sql.NullString) bool {
	table := u.database.Init()
	defer table.Close()
	tx, _ := table.Begin()
	defer func() {
		if err := recover(); err != nil {
			logging.SharedInstance().MethodInfo("list", "Save", true).Error("unexpected error")
			tx.Rollback()
		}
	}()

	result, err := tx.Exec("insert into lists (project_id, user_id, title, color, list_option_id, created_at) values (?, ?, ?, ?, ?, now());", u.ProjectID, u.UserID, u.Title, u.Color, u.ListOptionID)
	if err != nil {
		logging.SharedInstance().MethodInfo("list", "Save", true).Errorf("list save error: %v", err)
		tx.Rollback()
		return false
	}

	if OauthToken != nil && OauthToken.Valid && repo != nil {
		token := OauthToken.String
		label, err := hub.CheckLabelPresent(token, repo, &u.Title.String)
		if err != nil {
			tx.Rollback()
			logging.SharedInstance().MethodInfo("list", "Save", true).Errorf("check label error: %v", err)
			return false
		} else if label == nil {
			label, err = hub.CreateGithubLabel(token, repo, &u.Title.String, &u.Color.String)
			if err != nil {
				logging.SharedInstance().MethodInfo("list", "Save", true).Error("github label create failed")
				tx.Rollback()
				return false
			}
		} else {
			// createしようとしたときに存在している場合，それはあまり気にしなくて良い．むしろこれで同等の状態になる
			logging.SharedInstance().MethodInfo("list", "Save").Info("github label already exist")
		}
	}
	tx.Commit()
	u.ID, _ = result.LastInsertId()
	return true
}

func (u *ListStruct) Update(repo *repository.RepositoryStruct, OauthToken *sql.NullString, title *string, color *string, action *string) bool {
	table := u.database.Init()
	defer table.Close()

	// 初期リストに関しては一切編集を許可しない
	// 色は変えられても良いが，titleとactionは変えられては困る
	// 第一段階では色も含めてすべて固定とする
	if u.IsInitList() {
		logging.SharedInstance().MethodInfo("list", "Update", true).Error("cannot update initial list")
		return false
	}

	var listOptionID sql.NullInt64
	listOption := list_option.FindByAction(*action)
	if listOption == nil {
		logging.SharedInstance().MethodInfo("list", "Update").Debug("cannot find list_options, set null to list_option_id")
	} else {
		listOptionID.Int64 = listOption.ID
		listOptionID.Valid = true
	}

	tx, _ := table.Begin()
	defer func() {
		if err := recover(); err != nil {
			logging.SharedInstance().MethodInfo("list", "Update", true).Error("unexpected error")
			tx.Rollback()
		}
	}()

	_, err := tx.Exec("update lists set title = ?, color = ?, list_option_id = ? where id = ?;", *title, *color, listOptionID, u.ID)
	if err != nil {
		logging.SharedInstance().MethodInfo("list", "Update", true).Errorf("list update error: %v", err)
		tx.Rollback()
		return false
	}

	if OauthToken != nil && OauthToken.Valid && repo != nil {
		token := OauthToken.String
		// 編集前のラベルがそもそも存在しているかどうかを確認する
		existLabel, err := hub.CheckLabelPresent(token, repo, &u.Title.String)
		if err != nil {
			tx.Rollback()
			logging.SharedInstance().MethodInfo("list", "Update", true).Errorf("check label error: %v", err)
			return false
		} else if existLabel == nil {
			// editの場合ここに入る可能性はほとんどない
			// 編集前のラベルが存在しなければ新しく作るのと同義
			// もし存在していた場合は，エラーにしたい
			// あくまでgithub側のデータを正としたい．そしてgithub側からfasciaに同期をかけるのはここの責務ではない．
			// そのため，ここは素直にエラーにして，同期処理側をしっかり作りこむべき
			_, err := hub.CreateGithubLabel(token, repo, title, color)
			if err != nil {
				tx.Rollback()
				return false
			}
		} else {
			_, err := hub.UpdateGithubLabel(token, repo, &u.Title.String, title, color)
			if err != nil {
				tx.Rollback()
				return false
			}
		}
	}

	tx.Commit()
	u.Title = sql.NullString{String: *title, Valid: true}
	u.Color = sql.NullString{String: *color, Valid: true}
	u.ListOptionID = listOptionID
	return true
}

func (u *ListStruct) Tasks() []*task.TaskStruct {
	table := u.database.Init()
	defer table.Close()

	rows, _ := table.Query("select id, list_id, project_id, user_id, issue_number, title, description, pull_request, html_url from tasks where list_id = ? order by display_index;", u.ID)
	var slice []*task.TaskStruct
	for rows.Next() {
		var id, listID, userID, projectID int64
		var title, description string
		var issueNumber sql.NullInt64
		var pullRequest bool
		var htmlURL sql.NullString
		err := rows.Scan(&id, &listID, &projectID, &userID, &issueNumber, &title, &description, &pullRequest, &htmlURL)
		if err != nil {
			logging.SharedInstance().MethodInfo("List", "Tasks", true).Panic(err)
		}
		if listID == u.ID {
			l := task.NewTask(id, listID, projectID, userID, issueNumber, title, description, pullRequest, htmlURL)
			slice = append(slice, l)
		}
	}
	return slice
}

func (u *ListStruct) IsInitList() bool {
	for _, elem := range config.Element("init_list").(map[interface{}]interface{}) {
		if u.Title.String == elem.(string) {
			return true
		}
	}
	return false
}

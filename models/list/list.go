package list

import (
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/models/db"
	"github.com/h3poteto/fascia/models/list_option"
	"github.com/h3poteto/fascia/models/repository"
	"github.com/h3poteto/fascia/models/task"
	"github.com/h3poteto/fascia/modules/hub"
	"github.com/h3poteto/fascia/modules/logging"

	"database/sql"
	"runtime"

	"github.com/pkg/errors"
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
	IsHidden     bool
	database     *sql.DB
}

func NewList(id int64, projectID int64, userID int64, title string, color string, optionID sql.NullInt64, isHidden bool) *ListStruct {
	if projectID == 0 {
		return nil
	}
	nullTitle := sql.NullString{String: title, Valid: true}
	nullColor := sql.NullString{String: color, Valid: true}

	list := &ListStruct{ID: id, ProjectID: projectID, UserID: userID, Title: nullTitle, Color: nullColor, ListOptionID: optionID, IsHidden: isHidden}
	list.Initialize()
	return list
}

func FindList(projectID int64, listID int64) (*ListStruct, error) {
	database := db.SharedInstance().Connection
	var id, userID int64
	var title, color sql.NullString
	var optionID sql.NullInt64
	var isHidden bool
	rows, err := database.Query("select id, project_id, user_id, title, color, list_option_id, is_hidden from lists where id = ? AND project_id = ?;", listID, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	for rows.Next() {
		err = rows.Scan(&id, &projectID, &userID, &title, &color, &optionID, &isHidden)
		if err != nil {
			return nil, errors.Wrap(err, "sql scan error")
		}
	}
	if id != listID {
		return nil, errors.New("cannot find list or project did not contain list")
	}
	list := NewList(id, projectID, userID, title.String, color.String, optionID, isHidden)
	return list, nil

}

func (u *ListStruct) Initialize() {
	u.database = db.SharedInstance().Connection
}

func (u *ListStruct) Save(repo *repository.RepositoryStruct, OauthToken *sql.NullString) (e error) {
	tx, _ := u.database.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			switch ty := err.(type) {
			case runtime.Error:
				e = errors.Wrap(ty, "runtime error")
			case string:
				e = errors.New(err.(string))
			default:
				e = errors.New("unexpected error")
			}
		}
	}()

	result, err := tx.Exec("insert into lists (project_id, user_id, title, color, list_option_id, is_hidden, created_at) values (?, ?, ?, ?, ?, ?, now());", u.ProjectID, u.UserID, u.Title, u.Color, u.ListOptionID, u.IsHidden)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "sql execute error")
	}

	if OauthToken != nil && OauthToken.Valid && repo != nil {
		token := OauthToken.String
		label, err := hub.CheckLabelPresent(token, repo, &u.Title.String)
		if err != nil {
			tx.Rollback()
			return err
		} else if label == nil {
			label, err = hub.CreateGithubLabel(token, repo, &u.Title.String, &u.Color.String)
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// 色だけはこちら指定のものに変更したい
			_, err := hub.UpdateGithubLabel(token, repo, &u.Title.String, &u.Title.String, &u.Color.String)
			if err != nil {
				tx.Rollback()
				return err
			}
			logging.SharedInstance().MethodInfo("list", "Save").Info("github label already exist")
		}
	}
	tx.Commit()
	u.ID, _ = result.LastInsertId()
	return nil
}

// Update update and save list in database
func (u *ListStruct) Update(repo *repository.RepositoryStruct, OauthToken *sql.NullString, title *string, color *string, optionID *int64) (e error) {
	// 初期リストに関しては一切編集を許可しない
	// 色は変えられても良いが，titleとactionは変えられては困る
	// 第一段階では色も含めてすべて固定とする
	if u.IsInitList() {
		return errors.New("cannot update initial list")
	}

	var listOptionID sql.NullInt64
	listOption, err := list_option.FindByID(sql.NullInt64{Int64: *optionID, Valid: true})
	if err != nil {
		// list_optionはnullでも構わない
		// nullの場合は特にactionが発生しないだけ
		logging.SharedInstance().MethodInfo("list", "Update").Debugf("cannot find list_options, set null to list_option_id: %v", err)
	} else {
		listOptionID.Int64 = listOption.ID
		listOptionID.Valid = true
	}

	tx, _ := u.database.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			switch ty := err.(type) {
			case runtime.Error:
				e = errors.Wrap(ty, "runtime error")
			case string:
				e = errors.New(err.(string))
			default:
				e = errors.New("unexpected error")
			}
		}
	}()

	_, err = tx.Exec("update lists set title = ?, color = ?, list_option_id = ?, is_hidden = ? where id = ?;", *title, *color, listOptionID, u.IsHidden, u.ID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "sql execute error")
	}

	if OauthToken != nil && OauthToken.Valid && repo != nil {
		token := OauthToken.String
		// 編集前のラベルがそもそも存在しているかどうかを確認する
		existLabel, err := hub.CheckLabelPresent(token, repo, &u.Title.String)
		if err != nil {
			tx.Rollback()
			return err
		} else if existLabel == nil {
			// editの場合ここに入る可能性はほとんどない
			// 編集前のラベルが存在しなければ新しく作るのと同義
			// もし存在していた場合は，エラーにしたい
			// あくまでgithub側のデータを正としたい．そしてgithub側からfasciaに同期をかけるのはここの責務ではない．
			// そのため，ここは素直にエラーにして，同期処理側をしっかり作りこむべき
			_, err := hub.CreateGithubLabel(token, repo, title, color)
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			_, err := hub.UpdateGithubLabel(token, repo, &u.Title.String, title, color)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()
	u.Title = sql.NullString{String: *title, Valid: true}
	u.Color = sql.NullString{String: *color, Valid: true}
	u.ListOptionID = listOptionID
	return nil
}

// UpdateColor sync list color to github
func (u *ListStruct) UpdateColor() error {
	_, err := u.database.Exec("update lists set color = ? where id = ?;", u.Color.String, u.ID)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	return nil
}

// Tasks list up related a list
func (u *ListStruct) Tasks() ([]*task.TaskStruct, error) {
	var slice []*task.TaskStruct
	rows, err := u.database.Query("select id, list_id, project_id, user_id, issue_number, title, description, pull_request, html_url from tasks where list_id = ? order by display_index;", u.ID)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}

	for rows.Next() {
		var id, listID, userID, projectID int64
		var title, description string
		var issueNumber sql.NullInt64
		var pullRequest bool
		var htmlURL sql.NullString
		err := rows.Scan(&id, &listID, &projectID, &userID, &issueNumber, &title, &description, &pullRequest, &htmlURL)
		if err != nil {
			return nil, errors.Wrap(err, "sql select error")
		}
		if listID == u.ID {
			l := task.NewTask(id, listID, projectID, userID, issueNumber, title, description, pullRequest, htmlURL)
			slice = append(slice, l)
		}
	}
	return slice, nil
}

func (u *ListStruct) IsInitList() bool {
	for _, elem := range config.Element("init_list").(map[interface{}]interface{}) {
		if u.Title.String == elem.(string) {
			return true
		}
	}
	return false
}

// Hide can hide a list, it change is_hidden field
func (u *ListStruct) Hide() error {
	_, err := u.database.Exec("update lists set is_hidden = true where id = ?;", u.ID)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	u.IsHidden = true
	return nil
}

// Display can display a list, it change is_hidden filed
func (u *ListStruct) Display() error {
	_, err := u.database.Exec("update lists set is_hidden = false where id = ?;", u.ID)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	u.IsHidden = false
	return nil
}

// HasCloseAction check a list has close list_option
func (u *ListStruct) HasCloseAction() (bool, error) {
	option, err := list_option.FindByID(u.ListOptionID)
	if err != nil {
		return false, err
	}
	return option.CloseAction(), nil
}

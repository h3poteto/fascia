package list

import (
	"../../modules/hub"
	"../db"
	"../repository"
	"../task"
	"database/sql"
	"fmt"
)

type List interface {
	Save() bool
}

type ListStruct struct {
	Id        int64
	ProjectId int64
	UserId    int64
	Title     sql.NullString
	ListTasks []*task.TaskStruct
	Color     sql.NullString
	database  db.DB
}

func NewList(id int64, projectID int64, userID int64, title string, color string) *ListStruct {
	if projectID == 0 {
		return nil
	}
	nullTitle := sql.NullString{String: title, Valid: true}
	nullColor := sql.NullString{String: color, Valid: true}
	list := &ListStruct{Id: id, ProjectId: projectID, UserId: userID, Title: nullTitle, Color: nullColor}
	list.Initialize()
	return list
}

func FindList(projectID int64, listID int64) *ListStruct {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	var id, projectId, userId int64
	var title, color sql.NullString
	rows, _ := table.Query("select id, project_id, user_id, title, color from lists where id = ? AND project_id = ?;", listID, projectID)
	for rows.Next() {
		err := rows.Scan(&id, &projectId, &userId, &title, &color)
		if err != nil {
			panic(err.Error())
		}
	}
	if id != listID {
		fmt.Printf("cannot find list or project did not contain list: %v\n", listID)
		return nil
	} else {
		list := NewList(id, projectId, userId, title.String, color.String)
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
			fmt.Printf("panic: %+v\n", err)
			tx.Rollback()
		}
	}()

	result, err := tx.Exec("insert into lists (project_id, user_id, title, color, created_at) values (?, ?, ?, ?, now());", u.ProjectId, u.UserId, u.Title, u.Color)
	if err != nil {
		fmt.Printf("list save error: %+v\n", err)
		tx.Rollback()
		return false
	}

	if OauthToken != nil && OauthToken.Valid && repo != nil {
		token := OauthToken.String
		label, err := hub.CheckLabelPresent(token, repo, &u.Title.String)
		if err != nil {
			tx.Rollback()
			fmt.Printf("check label error: %+v\n", err)
			return false
		} else if label == nil {
			label, err = hub.CreateGithubLabel(token, repo, &u.Title.String, &u.Color.String)
			if err != nil {
				fmt.Printf("github label create failed\n")
				tx.Rollback()
				return false
			}
		} else {
			// createしようとしたときに存在している場合，それはあまり気にしなくて良い．むしろこれで同等の状態になる
			fmt.Printf("github label already exist\n")
		}
	}
	tx.Commit()
	u.Id, _ = result.LastInsertId()
	return true
}

func (u *ListStruct) Update(repo *repository.RepositoryStruct, OauthToken *sql.NullString, title *string, color *string) bool {
	table := u.database.Init()
	defer table.Close()
	tx, _ := table.Begin()
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic: %+v\n", err)
			tx.Rollback()
		}
	}()

	_, err := tx.Exec("update lists set title = ?, color = ? where id = ?;", *title, *color, u.Id)
	if err != nil {
		tx.Rollback()
		return false
	}

	if OauthToken != nil && OauthToken.Valid && repo != nil {
		token := OauthToken.String
		fmt.Printf("repository: %+v\n", repo)
		// 編集前のラベルがそもそも存在しているかどうかを確認する
		existLabel, err := hub.CheckLabelPresent(token, repo, &u.Title.String)
		fmt.Printf("find label: %+v\n", existLabel)
		if err != nil {
			tx.Rollback()
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
	return true
}

func (u *ListStruct) Tasks() []*task.TaskStruct {
	table := u.database.Init()
	defer table.Close()

	rows, _ := table.Query("select id, list_id, user_id, issue_number, title from tasks where list_id = ? order by display_index;", u.Id)
	var slice []*task.TaskStruct
	for rows.Next() {
		var id, listID, userID int64
		var title sql.NullString
		var issueNumber sql.NullInt64
		err := rows.Scan(&id, &listID, &userID, &issueNumber, &title)
		if err != nil {
			panic(err.Error())
		}
		if listID == u.Id && title.Valid {
			l := task.NewTask(id, listID, userID, issueNumber, title.String)
			slice = append(slice, l)
		}
	}
	return slice
}

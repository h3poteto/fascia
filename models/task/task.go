package task

import (
	"../../modules/hub"
	"../db"
	"../repository"
	"database/sql"
	"fmt"
)

type Task interface {
	Save() bool
}

type TaskStruct struct {
	Id       int64
	ListId   int64
	UserId   int64
	Title    sql.NullString
	database db.DB
}

func NewTask(id int64, listID int64, userID int64, title string) *TaskStruct {
	if listID == 0 {
		return nil
	}
	nullTitle := sql.NullString{String: title, Valid: true}
	task := &TaskStruct{Id: id, ListId: listID, UserId: userID, Title: nullTitle}
	task.Initialize()
	return task
}

func FindTask(listID int64, taskID int64) *TaskStruct {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	var id, listId, userId int64
	var title string
	rows, _ := table.Query("select id, list_id, user_id, title from tasks where id = ? AND list_id = ?;", taskID, listID)
	for rows.Next() {
		err := rows.Scan(&id, &listId, &userId, &title)
		if err != nil {
			panic(err.Error())
		}
	}
	if id != taskID {
		fmt.Printf("cannot find task or list did not contain task: %v\n", taskID)
		return nil
	} else {
		task := NewTask(id, listId, userId, title)
		return task
	}
}

func (u *TaskStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}

func (u *TaskStruct) Save(repo *repository.RepositoryStruct, OauthToken *sql.NullString) bool {
	table := u.database.Init()
	defer table.Close()
	transaction, _ := table.Begin()
	defer func() {
		if err := recover(); err != nil {
			transaction.Rollback()
		}
	}()

	// display_indexを自動挿入する
	count := 0
	err := transaction.QueryRow("SELECT COUNT(id) FROM tasks WHERE list_id = ?;", u.ListId).Scan(&count)
	result, err := transaction.Exec("insert into tasks (list_id, title, display_index, created_at) values (?, ?, ?, now());", u.ListId, u.Title, count+1)
	if err != nil {
		transaction.Rollback()
		return false
	}

	if OauthToken != nil && OauthToken.Valid && repo != nil {
		token := OauthToken.String
		label := hub.CheckLabelPresent(u.ListId, token, repo)
		// もしラベルがなかった場合は作っておく
		// 色が違っていてもアップデートは不要，それは編集でやってくれ
		if label == nil {
			label = hub.CreateGithubLabel(u.ListId, token, repo)
			if label == nil {
				transaction.Rollback()
				return false
			}
		}
		// issueを作る
		issue := hub.CreateGithubIssue(u.Id, token, repo, []string{*label.Name})
		if issue == nil {
			transaction.Rollback()
			return false
		}
	}

	err = transaction.Commit()
	if err != nil {
		transaction.Rollback()
		return false
	}
	u.Id, _ = result.LastInsertId()
	return true
}

// lastに追加する場合にはprevToTaskIdをnullで渡す
func (u *TaskStruct) ChangeList(listId int64, prevToTaskId *int64) bool {
	table := u.database.Init()
	defer table.Close()
	transaction, _ := table.Begin()
	defer func() {
		// panicがおきたらロールバック
		if err := recover(); err != nil {
			fmt.Printf("err: %+v\n", err)
			transaction.Rollback()
		}
	}()

	var prevToTaskIndex int
	if prevToTaskId != nil {
		err := transaction.QueryRow("select display_index from tasks where id = ?;", *prevToTaskId).Scan(&prevToTaskIndex)
		if err != nil {
			panic(err.Error())
		}
		// 先に後ろにいる奴らを押し出しておかないとprevToTaskIndexのg位置が開かない
		// prevToTaskIndex = nilのときは，末尾挿入なので払い出しは不要
		_, err = transaction.Exec("update tasks set display_index = display_index + 1 where id in (select id from (select id from tasks where list_id = ? and display_index >= ?) as tmp);", listId, prevToTaskIndex)
		if err != nil {
			panic(err.Error())
		}
	} else {
		// 本当は連番のはずだからカウントすればいいんだけど，念の為ラストのindex+1を取る
		// list内のタスクが空だった場合のためにnilが帰ってくることを許容する
		var index interface{}
		err := transaction.QueryRow("select max(display_index) from tasks where list_id = ?;", listId).Scan(&index)
		if err != nil {
			panic(err.Error())
		}
		if index == nil {
			prevToTaskIndex = 1
		} else {
			prevToTaskIndex = int(index.(int64)) + 1
		}
	}

	_, err := transaction.Exec("update tasks set list_id = ?, display_index = ? where id = ?;", listId, prevToTaskIndex, u.Id)
	if err != nil {
		panic(err.Error())
	}
	err = transaction.Commit()
	if err != nil {
		panic(err.Error())
	}
	u.ListId = listId
	return true
}

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
	Id          int64
	ListId      int64
	UserId      int64
	IssueNumber sql.NullInt64
	Title       sql.NullString
	database    db.DB
}

func NewTask(id int64, listID int64, userID int64, issueNumber sql.NullInt64, title string) *TaskStruct {
	if listID == 0 {
		return nil
	}
	nullTitle := sql.NullString{String: title, Valid: true}
	task := &TaskStruct{Id: id, ListId: listID, UserId: userID, IssueNumber: issueNumber, Title: nullTitle}
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
	var issueNumber sql.NullInt64
	err := table.QueryRow("select id, list_id, user_id, issue_number, title from tasks where id = ? AND list_id = ?;", taskID, listID).Scan(&id, &listId, &userId, &issueNumber, &title)
	if err != nil {
		panic(err.Error())
	}
	if id != taskID {
		fmt.Printf("cannot find task or list did not contain task: %v\n", taskID)
		return nil
	} else {
		task := NewTask(id, listId, userId, issueNumber, title)
		return task
	}
}

func FindByIssueNumber(issueNumber int) *TaskStruct {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	var id, listId, userId int64
	var title string
	var number sql.NullInt64
	err := table.QueryRow("select id, list_id, user_id, issue_number, title from tasks where issue_number = ?;", issueNumber).Scan(&id, &listId, &userId, &number, &title)
	if err != nil {
		panic(err.Error())
	}
	if !number.Valid || number.Int64 != int64(issueNumber) {
		fmt.Printf("cannot find task issue number: %v\n", issueNumber)
		return nil
	} else {
		task := NewTask(id, listId, userId, number, title)
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
			fmt.Printf("err: %+v\n", err)
			transaction.Rollback()
		}
	}()

	// display_indexを自動挿入する
	count := 0
	err := transaction.QueryRow("SELECT COUNT(id) FROM tasks WHERE list_id = ?;", u.ListId).Scan(&count)
	result, err := transaction.Exec("insert into tasks (list_id, user_id, title, display_index, created_at) values (?, ?, ?, ?, now());", u.ListId, u.UserId, u.Title, count+1)
	if err != nil {
		fmt.Printf("insert task error: %+v\n", err)
		transaction.Rollback()
		return false
	}
	var listTitle, listColor sql.NullString
	err = transaction.QueryRow("select title, color from lists where id = ?;", u.ListId).Scan(&listTitle, &listColor)
	if err != nil {
		fmt.Printf("select list error: %+v\n", err)
		transaction.Rollback()
		return false
	}
	if OauthToken != nil && OauthToken.Valid && repo != nil {
		token := OauthToken.String
		label := hub.CheckLabelPresent(token, repo, &listTitle.String)
		// もしラベルがなかった場合は作っておく
		// 色が違っていてもアップデートは不要，それは編集でやってくれ
		if label == nil {
			label = hub.CreateGithubLabel(token, repo, &listTitle.String, &listColor.String)
			if label == nil {
				transaction.Rollback()
				return false
			}
		}
		// issueを作る
		issue := hub.CreateGithubIssue(token, repo, []string{*label.Name}, &u.Title.String)
		if issue == nil {
			fmt.Printf("issue create failed:%+v\n", issue)
			transaction.Rollback()
			return false
		}
		currentId, _ := result.LastInsertId()
		_, err = transaction.Exec("update tasks set issue_number = ? where id = ?;", *issue.Number, currentId)
		fmt.Printf("issue number update\n")
		if err != nil {
			// TODO: そもそもこのときはissueを削除しなければいけないのでは？
			fmt.Printf("issue_number update error: %+v\n", err)
			transaction.Rollback()
			return false
		}
		u.IssueNumber = sql.NullInt64{Int64: int64(*issue.Number), Valid: true}
	}

	err = transaction.Commit()
	if err != nil {
		fmt.Printf("commit error:%+v\n", err)
		transaction.Rollback()
		return false
	}
	u.Id, _ = result.LastInsertId()
	return true
}

func (u *TaskStruct) Update(repo *repository.RepositoryStruct, OauthToken *sql.NullString) bool {
	table := u.database.Init()
	defer table.Close()

	_, err := table.Exec("update tasks set list_id = ?, issue_number = ?, title = ? where id = ?;", u.ListId, u.IssueNumber, u.Title, u.Id)
	if err != nil {
		fmt.Printf("update error: %+v\n", err)
		return false
	}
	return true
}

// lastに追加する場合にはprevToTaskIdをnullで渡す
func (u *TaskStruct) ChangeList(listId int64, prevToTaskId *int64, repo *repository.RepositoryStruct, OauthToken *sql.NullString) bool {
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

	// リストを移動させるのか同リスト内の並び替えなのかどうかを見て，並び替えならgithub同期したくない
	var isReorder bool
	if listId == u.ListId {
		isReorder = true
	} else {
		isReorder = false
	}

	var prevToTaskIndex int
	if prevToTaskId != nil {
		err := transaction.QueryRow("select display_index from tasks where id = ?;", *prevToTaskId).Scan(&prevToTaskIndex)
		if err != nil {
			fmt.Printf("select display index error: %+v\n", err)
			transaction.Rollback()
			return false
		}
		// 先に後ろにいる奴らを押し出しておかないとprevToTaskIndexのg位置が開かない
		// prevToTaskIndex = nilのときは，末尾挿入なので払い出しは不要
		_, err = transaction.Exec("update tasks set display_index = display_index + 1 where id in (select id from (select id from tasks where list_id = ? and display_index >= ?) as tmp);", listId, prevToTaskIndex)
		if err != nil {
			fmt.Printf("update display index error:%+v\n", err)
			transaction.Rollback()
			return false
		}
	} else {
		// 本当は連番のはずだからカウントすればいいんだけど，念の為ラストのindex+1を取る
		// list内のタスクが空だった場合のためにnilが帰ってくることを許容する
		var index interface{}
		err := transaction.QueryRow("select max(display_index) from tasks where list_id = ?;", listId).Scan(&index)
		if err != nil {
			fmt.Printf("select mas display index error:%+v\n", err)
			transaction.Rollback()
			return false
		}
		if index == nil {
			prevToTaskIndex = 1
		} else {
			prevToTaskIndex = int(index.(int64)) + 1
		}
	}

	_, err := transaction.Exec("update tasks set list_id = ?, display_index = ? where id = ?;", listId, prevToTaskIndex, u.Id)
	if err != nil {
		fmt.Printf("update task error:%+v\n", err)
		transaction.Rollback()
		return false
	}

	if !isReorder && OauthToken != nil && OauthToken.Valid && repo != nil && u.IssueNumber.Valid {
		token := OauthToken.String
		var listTitle, listColor sql.NullString
		err = transaction.QueryRow("select title, color from lists where id = ?;", listId).Scan(&listTitle, &listColor)
		if err != nil {
			fmt.Printf("select list error:%+v\n", err)
			transaction.Rollback()
			return false
		}
		label := hub.CheckLabelPresent(token, repo, &listTitle.String)
		if label == nil {
			// 移動先がない場合はつくろう
			label = hub.CreateGithubLabel(token, repo, &listTitle.String, &listColor.String)
			if label == nil {
				transaction.Rollback()
				return false
			}
		}
		// issueを移動
		if !hub.ReplaceLabelsForIssue(token, repo, u.IssueNumber.Int64, []string{*label.Name}) {
			transaction.Rollback()
			return false
		}
	}

	err = transaction.Commit()
	if err != nil {
		panic(err.Error())
	}
	u.ListId = listId
	return true
}

package task

import (
	"../../config"
	"../../modules/hub"
	"../../modules/logging"
	"../db"
	"../list_option"
	"../repository"
	"database/sql"
	"errors"
)

type Task interface {
	Save() bool
}

type TaskStruct struct {
	Id          int64
	ListId      int64
	UserId      int64
	IssueNumber sql.NullInt64
	Title       string
	Description string
	database    db.DB
}

func NewTask(id int64, listID int64, userID int64, issueNumber sql.NullInt64, title string, description string) *TaskStruct {
	task := &TaskStruct{Id: id, ListId: listID, UserId: userID, IssueNumber: issueNumber, Title: title, Description: description}
	task.Initialize()
	return task
}

func FindTask(listID int64, taskID int64) (*TaskStruct, error) {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	var id, listId, userId int64
	var title, description string
	var issueNumber sql.NullInt64
	err := table.QueryRow("select id, list_id, user_id, issue_number, title, description from tasks where id = ? AND list_id = ?;", taskID, listID).Scan(&id, &listId, &userId, &issueNumber, &title, &description)
	if err != nil {
		logging.SharedInstance().MethodInfo("task", "FindTask").Errorf("cannot find task: %v", err)
		return nil, err
	}
	if id != taskID {
		logging.SharedInstance().MethodInfo("task", "FindTask").Errorf("cannot find task or list did not contain task: %v", taskID)
		return nil, errors.New("cannot find task or list did not contain task")
	}
	task := NewTask(id, listId, userId, issueNumber, title, description)
	return task, nil
}

func FindByIssueNumber(issueNumber int) (*TaskStruct, error) {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	var id, listId, userId int64
	var title, description string
	var number sql.NullInt64
	err := table.QueryRow("select id, list_id, user_id, issue_number, title, description from tasks where issue_number = ?;", issueNumber).Scan(&id, &listId, &userId, &number, &title, &description)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	if !number.Valid || number.Int64 != int64(issueNumber) {
		logging.SharedInstance().MethodInfo("task", "FindByIssueNumber").Errorf("cannot find task issue number: %v", issueNumber)
		return nil, errors.New("task not found")
	} else {
		task := NewTask(id, listId, userId, number, title, description)
		return task, nil
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
			logging.SharedInstance().MethodInfo("task", "Save").Error("unexpected error")
			transaction.Rollback()
		}
	}()

	// display_indexを自動挿入する
	count := 0
	err := transaction.QueryRow("SELECT COUNT(id) FROM tasks WHERE list_id = ?;", u.ListId).Scan(&count)
	result, err := transaction.Exec("insert into tasks (list_id, user_id, issue_number, title, description, display_index, created_at) values (?, ?, ?, ?, ?, ?, now());", u.ListId, u.UserId, u.IssueNumber, u.Title, u.Description, count+1)
	if err != nil {
		logging.SharedInstance().MethodInfo("task", "Save").Errorf("insert task error: %v", err)
		transaction.Rollback()
		return false
	}
	if OauthToken != nil && OauthToken.Valid && repo != nil {
		var listTitle, listColor sql.NullString
		var listOptionId sql.NullInt64
		err = transaction.QueryRow("select title, color, list_option_id from lists where id = ?;", u.ListId).Scan(&listTitle, &listColor, &listOptionId)
		if err != nil {
			logging.SharedInstance().MethodInfo("task", "Save").Errorf("select list error: %v", err)
			transaction.Rollback()
			return false
		}

		token := OauthToken.String
		label, err := hub.CheckLabelPresent(token, repo, &listTitle.String)
		// もしラベルがなかった場合は作っておく
		// 色が違っていてもアップデートは不要，それは編集でやってくれ
		if err != nil {
			transaction.Rollback()
			logging.SharedInstance().MethodInfo("task", "Save").Errorf("check label error: %v", err)
			return false
		} else if label == nil {
			label, err = hub.CreateGithubLabel(token, repo, &listTitle.String, &listColor.String)
			if err != nil {
				transaction.Rollback()
				logging.SharedInstance().MethodInfo("task", "Save").Errorf("create label error: %v", err)
				return false
			}
		}
		// issueを作る
		issue, err := hub.CreateGithubIssue(token, repo, []string{*label.Name}, &u.Title, &u.Description)
		if err != nil {
			logging.SharedInstance().MethodInfo("task", "Save").Errorf("issue create failed:%v", err)
			transaction.Rollback()
			return false
		}
		currentId, _ := result.LastInsertId()
		_, err = transaction.Exec("update tasks set issue_number = ? where id = ?;", *issue.Number, currentId)
		if err != nil {
			// TODO: そもそもこのときはissueを削除しなければいけないのでは？
			logging.SharedInstance().MethodInfo("task", "Save").Errorf("issue_number update error: %v", err)
			transaction.Rollback()
			return false
		}
		logging.SharedInstance().MethodInfo("task", "Save").Info("issue number is updated")
		u.IssueNumber = sql.NullInt64{Int64: int64(*issue.Number), Valid: true}
	}

	err = transaction.Commit()
	if err != nil {
		logging.SharedInstance().MethodInfo("task", "Save").Errorf("commit error:%v", err)
		transaction.Rollback()
		return false
	}
	u.Id, _ = result.LastInsertId()
	return true
}

func (u *TaskStruct) Update(repo *repository.RepositoryStruct, OauthToken *sql.NullString) bool {
	table := u.database.Init()
	defer table.Close()

	_, err := table.Exec("update tasks set list_id = ?, issue_number = ?, title = ?, description = ? where id = ?;", u.ListId, u.IssueNumber, u.Title, u.Description, u.Id)
	if err != nil {
		logging.SharedInstance().MethodInfo("task", "Update").Errorf("update error: %v", err)
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
		if err := recover(); err != nil {
			logging.SharedInstance().MethodInfo("task", "ChangeList").Error("unexpected error")
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
		// 途中に入れるパターン
		err := transaction.QueryRow("select display_index from tasks where id = ?;", *prevToTaskId).Scan(&prevToTaskIndex)
		if err != nil {
			logging.SharedInstance().MethodInfo("task", "ChangeList").Errorf("select display index error: %v", err)
			transaction.Rollback()
			return false
		}
		// 先に後ろにいる奴らを押し出しておかないとprevToTaskIndexのg位置が開かない
		// prevToTaskIndex = nilのときは，末尾挿入なので払い出しは不要
		_, err = transaction.Exec("update tasks set display_index = display_index + 1 where id in (select id from (select id from tasks where list_id = ? and display_index >= ?) as tmp);", listId, prevToTaskIndex)
		if err != nil {
			logging.SharedInstance().MethodInfo("task", "ChangeList").Errorf("update display index error: %v", err)
			transaction.Rollback()
			return false
		}
	} else {
		// 最後尾に入れるパターン
		// 本当は連番のはずだからカウントすればいいんだけど，念の為ラストのindex+1を取る
		// list内のタスクが空だった場合のためにnilが帰ってくることを許容する
		var index interface{}
		err := transaction.QueryRow("select max(display_index) from tasks where list_id = ?;", listId).Scan(&index)
		if err != nil {
			// 該当するtaskが存在しないとき，indexにはnillが入るが，エラーにはならないので，ここのハンドリングには入らない
			logging.SharedInstance().MethodInfo("task", "ChangeList").Errorf("select max display index error:%v", err)
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
		logging.SharedInstance().MethodInfo("task", "ChangeList").Errorf("update task error:%v", err)
		transaction.Rollback()
		return false
	}

	// TODO: noneListの場合はlabelを外す処理
	if !isReorder && OauthToken != nil && OauthToken.Valid && repo != nil && u.IssueNumber.Valid {
		token := OauthToken.String
		var listTitle, listColor sql.NullString
		var listOptionId sql.NullInt64
		err = transaction.QueryRow("select title, color, list_option_id from lists where id = ?;", listId).Scan(&listTitle, &listColor, &listOptionId)
		if err != nil {
			logging.SharedInstance().MethodInfo("task", "ChangeList").Errorf("select list error:%v", err)
			transaction.Rollback()
			return false
		}

		var labelName []string
		if listTitle.String == config.Element("init_list").(map[interface{}]interface{})["none"].(string) {
			labelName = []string{}
		} else {
			label, err := hub.CheckLabelPresent(token, repo, &listTitle.String)
			if err != nil {
				logging.SharedInstance().MethodInfo("task", "ChangeList").Errorf("check label error: %v", err)
				transaction.Rollback()
				return false
			} else if label == nil {
				// 移動先がない場合はつくろう
				label, err = hub.CreateGithubLabel(token, repo, &listTitle.String, &listColor.String)
				if label == nil {
					logging.SharedInstance().MethodInfo("task", "ChangeList").Errorf("create label error: %v", err)
					transaction.Rollback()
					return false
				}
			}
			labelName = []string{*label.Name}
		}
		// list_option
		var issueAction *string
		listOption := list_option.FindById(listOptionId)
		if listOption != nil {
			issueAction = &listOption.Action
		}
		// issueを移動
		result, err := hub.EditGithubIssue(token, repo, u.IssueNumber.Int64, labelName, &u.Title, &u.Description, issueAction)
		if err != nil || !result {
			transaction.Rollback()
			return false
		}
	}

	err = transaction.Commit()
	if err != nil {
		logging.SharedInstance().MethodInfo("Task", "ChangeList").Panic(err)
	}
	u.ListId = listId
	return true
}

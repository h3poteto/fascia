package task

import (
	"../../config"
	"../../modules/hub"
	"../../modules/logging"
	"../db"
	"../list_option"
	"../repository"

	"database/sql"
	"runtime"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

type Task interface {
	Save() bool
}

type TaskStruct struct {
	ID          int64
	ListID      int64
	ProjectID   int64
	UserID      int64
	IssueNumber sql.NullInt64
	Title       string
	Description string
	PullRequest bool
	HTMLURL     sql.NullString
	database    db.DB
}

func NewTask(id int64, listID int64, projectID int64, userID int64, issueNumber sql.NullInt64, title string, description string, pullRequest bool, htmlURL sql.NullString) *TaskStruct {
	task := &TaskStruct{ID: id, ListID: listID, ProjectID: projectID, UserID: userID, IssueNumber: issueNumber, Title: title, Description: description, PullRequest: pullRequest, HTMLURL: htmlURL}
	task.Initialize()
	return task
}

func FindTask(listID int64, taskID int64) (*TaskStruct, error) {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	var id, userID, projectID int64
	var title, description string
	var issueNumber sql.NullInt64
	var pullRequest bool
	var htmlURL sql.NullString
	err := table.QueryRow("select id, list_id, project_id, user_id, issue_number, title, description, pull_request, html_url from tasks where id = ? AND list_id = ?;", taskID, listID).Scan(&id, &listID, &projectID, &userID, &issueNumber, &title, &description, &pullRequest, &htmlURL)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	if id != taskID {
		return nil, errors.New("cannot find task or list did not contain task")
	}
	task := NewTask(id, listID, projectID, userID, issueNumber, title, description, pullRequest, htmlURL)
	return task, nil
}

func FindByIssueNumber(projectID int64, issueNumber int) (*TaskStruct, error) {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	var id, listID, userID int64
	var title, description string
	var number sql.NullInt64
	var pullRequest bool
	var htmlURL sql.NullString
	err := table.QueryRow("select id, list_id, project_id, user_id, issue_number, title, description, pull_request, html_url from tasks where issue_number = ? and project_id = ?;", issueNumber, projectID).Scan(&id, &listID, &projectID, &userID, &number, &title, &description, &pullRequest, &htmlURL)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	if !number.Valid || number.Int64 != int64(issueNumber) {
		return nil, errors.New("task not found")
	}
	task := NewTask(id, listID, projectID, userID, number, title, description, pullRequest, htmlURL)
	return task, nil
}

func (u *TaskStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}

func (u *TaskStruct) Save(repo *repository.RepositoryStruct, OauthToken *sql.NullString) (e error) {
	table := u.database.Init()
	defer table.Close()
	transaction, _ := table.Begin()
	defer func() {
		if err := recover(); err != nil {
			transaction.Rollback()
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

	// display_indexを自動挿入する
	count := 0
	err := transaction.QueryRow("SELECT COUNT(id) FROM tasks WHERE list_id = ?;", u.ListID).Scan(&count)
	if err != nil {
		transaction.Rollback()
		return errors.Wrap(err, "sql select error")
	}
	result, err := transaction.Exec("insert into tasks (list_id, project_id, user_id, issue_number, title, description, pull_request, html_url, display_index, created_at) values (?,?,?, ?, ?, ?, ?, ?, ?, now());", u.ListID, u.ProjectID, u.UserID, u.IssueNumber, u.Title, u.Description, u.PullRequest, u.HTMLURL, count+1)
	if err != nil {
		transaction.Rollback()
		return errors.Wrap(err, "sql execute error")
	}
	currentID, _ := result.LastInsertId()

	// TODO: ここもsyncIssueを使いたい
	if OauthToken != nil && OauthToken.Valid && repo != nil {
		var listTitle, listColor sql.NullString
		var listOptionID sql.NullInt64
		err = transaction.QueryRow("select title, color, list_option_id from lists where id = ?;", u.ListID).Scan(&listTitle, &listColor, &listOptionID)
		if err != nil {
			transaction.Rollback()
			return errors.Wrap(err, "sql select error")
		}

		token := OauthToken.String

		labelName, err := u.syncLabel(listTitle.String, listColor.String, token, repo)
		if err != nil {
			transaction.Rollback()
			return err
		}

		// issueを作る
		issue, err := hub.CreateGithubIssue(token, repo, labelName, &u.Title, &u.Description)
		if err != nil {
			transaction.Rollback()
			return err
		}

		_, err = transaction.Exec("update tasks set issue_number = ?, pull_request = false, html_url = ? where id = ?;", *issue.Number, *issue.HTMLURL, currentID)
		if err != nil {
			// TODO: そもそもこのときはissueを削除しなければいけないのでは？
			// しかしissueの削除は不可能のはずで，どうするかを考えないといけないが，そもそもここの発生確率って今のところかなり低いはずなので，そこまで気にする必要はないのではないか？
			transaction.Rollback()
			return errors.Wrap(err, "sql execute error")
		}
		logging.SharedInstance().MethodInfo("task", "Save").Info("issue number is updated")
		u.IssueNumber = sql.NullInt64{Int64: int64(*issue.Number), Valid: true}
		u.HTMLURL = sql.NullString{String: *issue.HTMLURL, Valid: true}
	}

	err = transaction.Commit()
	if err != nil {
		transaction.Rollback()
		return errors.Wrap(err, "sql execute error")
	}
	u.ID, _ = result.LastInsertId()
	logging.SharedInstance().MethodInfo("task", "Save").Debugf("new task saved: %+v", u)
	return nil
}

// Update is update task in db and sync github issue
func (u *TaskStruct) Update(repo *repository.RepositoryStruct, OauthToken *sql.NullString) error {
	table := u.database.Init()
	defer table.Close()

	_, err := table.Exec("update tasks set list_id = ?, issue_number = ?, title = ?, description = ?, pull_request = ?, html_url = ? where id = ?;", u.ListID, u.IssueNumber, u.Title, u.Description, u.PullRequest, u.HTMLURL, u.ID)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	logging.SharedInstance().MethodInfo("task", "Update").Debugf("task updated: %+v", u)

	// github側へ同期
	if repo != nil && OauthToken != nil && OauthToken.Valid {
		go func() {
			_, err := u.syncIssue(repo, OauthToken)
			if err != nil {
				logging.SharedInstance().MethodInfo("task", "Update").Error(err)
				return
			}
			logging.SharedInstance().MethodInfo("task", "Update").Debugf("task synced to github: %+v", u)
			return
		}()
	}
	return nil
}

// lastに追加する場合にはprevToTaskIDをnullで渡す
func (u *TaskStruct) ChangeList(listID int64, prevToTaskID *int64, repo *repository.RepositoryStruct, OauthToken *sql.NullString) (e error) {
	table := u.database.Init()
	defer table.Close()
	transaction, _ := table.Begin()
	defer func() {
		if err := recover(); err != nil {
			transaction.Rollback()
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

	// リストを移動させるのか同リスト内の並び替えなのかどうかを見て，並び替えならgithub同期したくない
	var isReorder bool
	if listID == u.ListID {
		isReorder = true
	} else {
		isReorder = false
	}

	var prevToTaskIndex int
	if prevToTaskID != nil {
		// 途中に入れるパターン
		err := transaction.QueryRow("select display_index from tasks where id = ?;", *prevToTaskID).Scan(&prevToTaskIndex)
		if err != nil {
			transaction.Rollback()
			return errors.Wrap(err, "sql select error")
		}
		// 先に後ろにいる奴らを押し出しておかないとprevToTaskIndexのg位置が開かない
		// prevToTaskIndex = nilのときは，末尾挿入なので払い出しは不要
		_, err = transaction.Exec("update tasks set display_index = display_index + 1 where id in (select id from (select id from tasks where list_id = ? and display_index >= ?) as tmp);", listID, prevToTaskIndex)
		if err != nil {
			transaction.Rollback()
			return errors.Wrap(err, "sql execute error")
		}
	} else {
		// 最後尾に入れるパターン
		// 本当は連番のはずだからカウントすればいいんだけど，念の為ラストのindex+1を取る
		// list内のタスクが空だった場合のためにnilが帰ってくることを許容する
		var index interface{}
		err := transaction.QueryRow("select max(display_index) from tasks where list_id = ?;", listID).Scan(&index)
		if err != nil {
			// 該当するtaskが存在しないとき，indexにはnillが入るが，エラーにはならないので，ここのハンドリングには入らない
			transaction.Rollback()
			return errors.Wrap(err, "sql select error")
		}
		if index == nil {
			prevToTaskIndex = 1
		} else {
			prevToTaskIndex = int(index.(int64)) + 1
		}
	}

	_, err := transaction.Exec("update tasks set list_id = ?, display_index = ? where id = ?;", listID, prevToTaskIndex, u.ID)
	if err != nil {
		transaction.Rollback()
		return errors.Wrap(err, "sql execute error")
	}

	err = transaction.Commit()
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	u.ListID = listID

	if !isReorder && OauthToken != nil && OauthToken.Valid && repo != nil {
		go func() {
			_, err := u.syncIssue(repo, OauthToken)
			if err != nil {
				logging.SharedInstance().MethodInfo("task", "ChangeList").Error(err)
				return
			}
			logging.SharedInstance().MethodInfo("Task", "Update").Debugf("task synced to github: %+v", u)
			return
		}()
	}

	return nil
}

func (u *TaskStruct) syncLabel(listTitle string, listColor string, token string, repo *repository.RepositoryStruct) ([]string, error) {
	if listTitle == config.Element("init_list").(map[interface{}]interface{})["none"].(string) {
		return []string{}, nil
	}
	label, err := hub.CheckLabelPresent(token, repo, &listTitle)
	if err != nil {
		return nil, err
	} else if label == nil {
		// 対象のラベルがない場合には新規作成する
		label, err = hub.CreateGithubLabel(token, repo, &listTitle, &listColor)
		if label == nil {
			return nil, err
		}
	}
	return []string{*label.Name}, nil
}

func (u *TaskStruct) syncIssue(repo *repository.RepositoryStruct, OauthToken *sql.NullString) (*github.Issue, error) {
	table := u.database.Init()
	defer table.Close()

	var listTitle, listColor sql.NullString
	var listOptionID sql.NullInt64
	err := table.QueryRow("select title, color, list_option_id from lists where id = ?;", u.ListID).Scan(&listTitle, &listColor, &listOptionID)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}

	token := OauthToken.String

	labelName, err := u.syncLabel(listTitle.String, listColor.String, token, repo)
	if err != nil {
		return nil, err
	}

	var issue *github.Issue
	// issueを確認する
	if u.IssueNumber.Valid {
		issue, err = hub.GetGithubIssue(token, repo, int(u.IssueNumber.Int64))
		if err != nil {
			return nil, err
		}
	}
	// issueがない場合には作成する
	if issue == nil {
		issue, err = hub.CreateGithubIssue(token, repo, labelName, &u.Title, &u.Description)
		if err != nil {
			return nil, err
		}
	}
	// issueがある場合には更新する
	// list_optionに応じてopen, closeを更新する
	var issueAction *string
	listOption, err := list_option.FindByID(listOptionID)
	if err == nil {
		issueAction = &listOption.Action
	}

	result, err := hub.EditGithubIssue(token, repo, *issue.Number, labelName, &u.Title, &u.Description, issueAction)
	if err != nil {
		return nil, err
	}
	if !result {
		return nil, errors.New("unexpected error")
	}
	return issue, nil
}

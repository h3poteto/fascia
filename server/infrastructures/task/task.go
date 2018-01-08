package task

import (
	"github.com/h3poteto/fascia/lib/modules/database"
	"github.com/h3poteto/fascia/lib/modules/logging"

	"database/sql"

	"github.com/pkg/errors"
)

// Task has task record
type Task struct {
	ID          int64
	ListID      int64
	ProjectID   int64
	UserID      int64
	IssueNumber sql.NullInt64
	Title       string
	Description string
	PullRequest bool
	HTMLURL     sql.NullString
	db          *sql.DB
}

// New returns a task object
func New(id int64, listID int64, projectID int64, userID int64, issueNumber sql.NullInt64, title string, description string, pullRequest bool, htmlURL sql.NullString) *Task {
	task := &Task{ID: id, ListID: listID, ProjectID: projectID, UserID: userID, IssueNumber: issueNumber, Title: title, Description: description, PullRequest: pullRequest, HTMLURL: htmlURL}
	task.initialize()
	return task
}

// Find search a task according to id
func Find(listID int64, taskID int64) (*Task, error) {
	db := database.SharedInstance().Connection

	var id, userID, projectID int64
	var title, description string
	var issueNumber sql.NullInt64
	var pullRequest bool
	var htmlURL sql.NullString
	err := db.QueryRow("select id, list_id, project_id, user_id, issue_number, title, description, pull_request, html_url from tasks where id = ? AND list_id = ?;", taskID, listID).Scan(&id, &listID, &projectID, &userID, &issueNumber, &title, &description, &pullRequest, &htmlURL)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	if id != taskID {
		return nil, errors.New("cannot find task or list did not contain task")
	}
	task := New(id, listID, projectID, userID, issueNumber, title, description, pullRequest, htmlURL)
	return task, nil
}

// FindByIssueNumber search a task according to issue number in github
func FindByIssueNumber(projectID int64, issueNumber int) (*Task, error) {
	db := database.SharedInstance().Connection

	var id, listID, userID int64
	var title, description string
	var number sql.NullInt64
	var pullRequest bool
	var htmlURL sql.NullString
	err := db.QueryRow("select id, list_id, project_id, user_id, issue_number, title, description, pull_request, html_url from tasks where issue_number = ? and project_id = ?;", issueNumber, projectID).Scan(&id, &listID, &projectID, &userID, &number, &title, &description, &pullRequest, &htmlURL)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	if !number.Valid || number.Int64 != int64(issueNumber) {
		return nil, errors.New("task not found")
	}
	task := New(id, listID, projectID, userID, number, title, description, pullRequest, htmlURL)
	return task, nil
}

func Tasks(parentListID int64) ([]*Task, error) {
	db := database.SharedInstance().Connection

	var slice []*Task
	rows, err := db.Query("select id, list_id, project_id, user_id, issue_number, title, description, pull_request, html_url from tasks where list_id = ? order by display_index;", parentListID)
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
		if listID == parentListID {
			l := New(id, listID, projectID, userID, issueNumber, title, description, pullRequest, htmlURL)
			slice = append(slice, l)
		}
	}
	return slice, nil
}

func (t *Task) initialize() {
	t.db = database.SharedInstance().Connection
}

// Save save task model in database, and arrange order tasks
func (t *Task) Save() error {
	transaction, err := t.db.Begin()
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}

	// display_indexを自動挿入する
	count := 0
	err = transaction.QueryRow("SELECT COUNT(id) FROM tasks WHERE list_id = ?;", t.ListID).Scan(&count)
	if err != nil {
		transaction.Rollback()
		return errors.Wrap(err, "sql select error")
	}
	result, err := transaction.Exec("insert into tasks (list_id, project_id, user_id, issue_number, title, description, pull_request, html_url, display_index, created_at) values (?,?,?, ?, ?, ?, ?, ?, ?, now());", t.ListID, t.ProjectID, t.UserID, t.IssueNumber, t.Title, t.Description, t.PullRequest, t.HTMLURL, count+1)
	if err != nil {
		transaction.Rollback()
		return errors.Wrap(err, "sql execute error")
	}
	t.ID, _ = result.LastInsertId()
	err = transaction.Commit()
	if err != nil {
		transaction.Rollback()
		return errors.Wrap(err, "sql execute error")
	}
	logging.SharedInstance().MethodInfo("task", "Save").Debugf("new task saved: %+v", t)
	return nil
}

// Update is update task in database
func (t *Task) Update(listID int64, issueNumber sql.NullInt64, title, description string, pullRequest bool, htmlURL sql.NullString) error {
	_, err := t.db.Exec("update tasks set list_id = ?, issue_number = ?, title = ?, description = ?, pull_request = ?, html_url = ? where id = ?;", listID, issueNumber, title, description, pullRequest, htmlURL, t.ID)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	t.ListID = listID
	t.IssueNumber = issueNumber
	t.Title = title
	t.Description = description
	t.PullRequest = pullRequest
	t.HTMLURL = htmlURL

	logging.SharedInstance().MethodInfo("task", "Update").Debugf("task updated: %+v", t)

	return nil
}

// ChangeList change list which is belonged a task
// If add task in bottom, transmit null to prevToTaskID
func (t *Task) ChangeList(listID int64, prevToTaskID *int64) error {
	transaction, err := t.db.Begin()
	if err != nil {
		return errors.Wrap(err, "sql execute error")
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

	_, err = transaction.Exec("update tasks set list_id = ?, display_index = ? where id = ?;", listID, prevToTaskIndex, t.ID)
	if err != nil {
		transaction.Rollback()
		return errors.Wrap(err, "sql execute error")
	}

	err = transaction.Commit()
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	t.ListID = listID

	return nil
}

// Delete is delete a task in db
func (t *Task) Delete() error {
	if t.IssueNumber.Valid {
		return errors.New("cannot delete")
	}

	_, err := t.db.Exec("delete from tasks where id = ?;", t.ID)
	if err != nil {
		return errors.Wrap(err, "sql delelet error")
	}
	logging.SharedInstance().MethodInfo("task", "Delete").Infof("task deleted: %v", t.ID)
	t.ID = 0
	return nil
}

func (t *Task) List() (sql.NullString, sql.NullString, sql.NullInt64, error) {
	var listTitle, listColor sql.NullString
	var listOptionID sql.NullInt64
	err := t.db.QueryRow("select title, color, list_option_id from lists where id = ?;", t.ListID).Scan(&listTitle, &listColor, &listOptionID)
	if err != nil {
		return sql.NullString{}, sql.NullString{}, sql.NullInt64{}, errors.Wrap(err, "sql select error")
	}
	return listTitle, listColor, listOptionID, nil
}

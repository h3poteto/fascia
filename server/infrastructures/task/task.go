package task

import (
	"github.com/h3poteto/fascia/lib/modules/logging"

	"database/sql"

	"github.com/pkg/errors"
)

// Task has db connection.
type Task struct {
	db *sql.DB
}

// New returns a task object
func New(db *sql.DB) *Task {
	return &Task{
		db,
	}
}

// Find search a task according to id.
func (t *Task) Find(id int64) (int64, int64, int64, int64, sql.NullInt64, string, string, bool, sql.NullString, error) {
	var listID, userID, projectID int64
	var title, description string
	var issueNumber sql.NullInt64
	var pullRequest bool
	var htmlURL sql.NullString
	err := t.db.QueryRow("select id, list_id, project_id, user_id, issue_number, title, description, pull_request, html_url from tasks where id = ?;", id).Scan(&id, &listID, &projectID, &userID, &issueNumber, &title, &description, &pullRequest, &htmlURL)
	if err != nil {
		return 0, 0, 0, 0, sql.NullInt64{}, "", "", false, sql.NullString{}, errors.Wrap(err, "task repository")
	}
	return id, listID, projectID, userID, issueNumber, title, description, pullRequest, htmlURL, nil
}

// FindByIssueNumber search a task according to issue number in github
func (t *Task) FindByIssueNumber(projectID int64, issueNumber int) (int64, int64, int64, int64, sql.NullInt64, string, string, bool, sql.NullString, error) {
	var id, listID, userID int64
	var title, description string
	var number sql.NullInt64
	var pullRequest bool
	var htmlURL sql.NullString
	err := t.db.QueryRow("select id, list_id, project_id, user_id, issue_number, title, description, pull_request, html_url from tasks where issue_number = ? and project_id = ?;", issueNumber, projectID).Scan(&id, &listID, &projectID, &userID, &number, &title, &description, &pullRequest, &htmlURL)
	if err != nil {
		return 0, 0, 0, 0, sql.NullInt64{}, "", "", false, sql.NullString{}, errors.Wrap(err, "task repository")
	}
	if !number.Valid || number.Int64 != int64(issueNumber) {
		return 0, 0, 0, 0, sql.NullInt64{}, "", "", false, sql.NullString{}, errors.New("task not found")
	}
	return id, listID, projectID, userID, number, title, description, pullRequest, htmlURL, nil
}

// Tasks returns all tasks related a list.
func (t *Task) Tasks(parentListID int64) ([]map[string]interface{}, error) {
	result := []map[string]interface{}{}
	rows, err := t.db.Query("select id, list_id, project_id, user_id, issue_number, title, description, pull_request, html_url from tasks where list_id = ? order by display_index;", parentListID)
	if err != nil {
		return result, errors.Wrap(err, "task repository")
	}

	for rows.Next() {
		var id, listID, userID, projectID int64
		var title, description string
		var issueNumber sql.NullInt64
		var pullRequest bool
		var htmlURL sql.NullString
		err := rows.Scan(&id, &listID, &projectID, &userID, &issueNumber, &title, &description, &pullRequest, &htmlURL)
		if err != nil {
			return []map[string]interface{}{}, errors.Wrap(err, "task repository")
		}
		if listID == parentListID {
			l := map[string]interface{}{
				"id":          id,
				"listID":      listID,
				"projectID":   projectID,
				"userID":      userID,
				"issueNumber": issueNumber,
				"title":       title,
				"description": description,
				"pullRequest": pullRequest,
				"htmlURL":     htmlURL,
			}
			result = append(result, l)
		}
	}
	return result, nil
}

// NonIssueTasks returns all tasks related a list.
func (t *Task) NonIssueTasks(projectID, userID int64) ([]map[string]interface{}, error) {
	result := []map[string]interface{}{}
	rows, err := t.db.Query("SELECT id, list_id, project_id, user_id, issue_number, title, description, pull_request, html_url FROM tasks WHERE project_id = ? and user_id = ? and issue_number IS NULL;", projectID, userID)
	if err != nil {
		return result, errors.Wrap(err, "task repository")
	}

	for rows.Next() {
		var id, listID, userID, projectID int64
		var title, description string
		var issueNumber sql.NullInt64
		var pullRequest bool
		var htmlURL sql.NullString
		err := rows.Scan(&id, &listID, &projectID, &userID, &issueNumber, &title, &description, &pullRequest, &htmlURL)
		if err != nil {
			return []map[string]interface{}{}, errors.Wrap(err, "task repository")
		}
		l := map[string]interface{}{
			"id":          id,
			"listID":      listID,
			"projectID":   projectID,
			"userID":      userID,
			"issueNumber": issueNumber,
			"title":       title,
			"description": description,
			"pullRequest": pullRequest,
			"htmlURL":     htmlURL,
		}
		result = append(result, l)
	}
	return result, nil
}

// Create save task model in database, and arrange order tasks
func (t *Task) Create(listID int64, projectID int64, userID int64, issueNumber sql.NullInt64, title string, description string, pullRequest bool, htmlURL sql.NullString) (int64, error) {
	transaction, err := t.db.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "task repository")
	}

	// display_indexを自動挿入する
	count := 0
	err = transaction.QueryRow("SELECT COUNT(id) FROM tasks WHERE list_id = ?;", listID).Scan(&count)
	if err != nil {
		transaction.Rollback()
		return 0, errors.Wrap(err, "task repository")
	}
	result, err := transaction.Exec("insert into tasks (list_id, project_id, user_id, issue_number, title, description, pull_request, html_url, display_index, created_at) values (?,?,?, ?, ?, ?, ?, ?, ?, now());", listID, projectID, userID, issueNumber, title, description, pullRequest, htmlURL, count+1)
	if err != nil {
		transaction.Rollback()
		return 0, errors.Wrap(err, "task repository")
	}
	id, _ := result.LastInsertId()
	err = transaction.Commit()
	if err != nil {
		transaction.Rollback()
		return 0, errors.Wrap(err, "task repository")
	}
	logging.SharedInstance().MethodInfo("task", "Save").Debugf("new task saved: %+v", t)
	return id, nil
}

// Update is update task in database
func (t *Task) Update(id, listID, projectID, userID int64, issueNumber sql.NullInt64, title, description string, pullRequest bool, htmlURL sql.NullString) error {
	_, err := t.db.Exec("update tasks set list_id = ?, project_id = ?, user_id = ?, issue_number = ?, title = ?, description = ?, pull_request = ?, html_url = ? where id = ?;", listID, projectID, userID, issueNumber, title, description, pullRequest, htmlURL, id)
	if err != nil {
		return errors.Wrap(err, "task repository")
	}
	logging.SharedInstance().MethodInfo("task", "Update").Debugf("task updated: %+v", t)

	return nil
}

// ChangeList change list which is belonged a task
// If add task in bottom, transmit null to prevToTaskID
func (t *Task) ChangeList(id, listID int64, prevToTaskID *int64) error {
	transaction, err := t.db.Begin()
	if err != nil {
		return errors.Wrap(err, "task repository")
	}

	// TODO: ロジックをentityに移動させたい
	var prevToTaskIndex int
	if prevToTaskID != nil {
		// 途中に入れるパターン
		err := transaction.QueryRow("select display_index from tasks where id = ?;", *prevToTaskID).Scan(&prevToTaskIndex)
		if err != nil {
			transaction.Rollback()
			return errors.Wrap(err, "task repository")
		}
		// 先に後ろにいる奴らを押し出しておかないとprevToTaskIndexのg位置が開かない
		// prevToTaskIndex = nilのときは，末尾挿入なので払い出しは不要
		_, err = transaction.Exec("update tasks set display_index = display_index + 1 where id in (select id from (select id from tasks where list_id = ? and display_index >= ?) as tmp);", listID, prevToTaskIndex)
		if err != nil {
			transaction.Rollback()
			return errors.Wrap(err, "task repository")
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
			return errors.Wrap(err, "task repository")
		}
		if index == nil {
			prevToTaskIndex = 1
		} else {
			prevToTaskIndex = int(index.(int64)) + 1
		}
	}

	_, err = transaction.Exec("update tasks set list_id = ?, display_index = ? where id = ?;", listID, prevToTaskIndex, id)
	if err != nil {
		transaction.Rollback()
		return errors.Wrap(err, "task repository")
	}

	err = transaction.Commit()
	if err != nil {
		return errors.Wrap(err, "task repository")
	}
	return nil
}

// Delete is delete a task in db
func (t *Task) Delete(id int64) error {
	_, err := t.db.Exec("delete from tasks where id = ?;", id)
	if err != nil {
		return errors.Wrap(err, "task repository")
	}
	logging.SharedInstance().MethodInfo("task", "Delete").Infof("task deleted: %v", id)
	return nil
}

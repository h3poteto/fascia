package task

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/domains/task"

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
func (t *Task) Find(id int64) (*task.Task, error) {
	var listID, userID, projectID, displayIndex int64
	var title, description string
	var issueNumber sql.NullInt64
	var pullRequest bool
	var htmlURL sql.NullString
	err := t.db.QueryRow("select id, list_id, project_id, user_id, issue_number, title, description, pull_request, html_url, display_index from tasks where id = ?;", id).Scan(&id, &listID, &projectID, &userID, &issueNumber, &title, &description, &pullRequest, &htmlURL, &displayIndex)
	if err != nil {
		return nil, errors.Wrap(err, "task repository")
	}
	return &task.Task{
		ID:           id,
		ListID:       listID,
		ProjectID:    projectID,
		UserID:       userID,
		IssueNumber:  issueNumber,
		Title:        title,
		Description:  description,
		PullRequest:  pullRequest,
		HTMLURL:      htmlURL,
		DisplayIndex: displayIndex,
	}, nil
}

// FindByIssueNumber search a task according to issue number in github
func (t *Task) FindByIssueNumber(projectID int64, issueNumber int) (*task.Task, error) {
	var id, listID, userID, displayIndex int64
	var title, description string
	var number sql.NullInt64
	var pullRequest bool
	var htmlURL sql.NullString
	err := t.db.QueryRow("select id, list_id, project_id, user_id, issue_number, title, description, pull_request, html_url, display_index from tasks where issue_number = ? and project_id = ?;", issueNumber, projectID).Scan(&id, &listID, &projectID, &userID, &number, &title, &description, &pullRequest, &htmlURL, &displayIndex)
	if err != nil {
		return nil, errors.Wrap(err, "task repository")
	}
	if !number.Valid || number.Int64 != int64(issueNumber) {
		return nil, errors.New("task not found")
	}
	return &task.Task{
		ID:           id,
		ListID:       listID,
		ProjectID:    projectID,
		UserID:       userID,
		IssueNumber:  sql.NullInt64{Int64: int64(issueNumber), Valid: true},
		Title:        title,
		Description:  description,
		PullRequest:  pullRequest,
		HTMLURL:      htmlURL,
		DisplayIndex: displayIndex,
	}, nil
}

// Tasks returns all tasks related a list.
func (t *Task) Tasks(parentListID int64) ([]*task.Task, error) {
	result := []*task.Task{}
	rows, err := t.db.Query("select id, list_id, project_id, user_id, issue_number, title, description, pull_request, html_url, display_index from tasks where list_id = ? order by display_index;", parentListID)
	if err != nil {
		return nil, errors.Wrap(err, "task repository")
	}

	for rows.Next() {
		var id, listID, userID, projectID, displayIndex int64
		var title, description string
		var issueNumber sql.NullInt64
		var pullRequest bool
		var htmlURL sql.NullString
		err := rows.Scan(&id, &listID, &projectID, &userID, &issueNumber, &title, &description, &pullRequest, &htmlURL, &displayIndex)
		if err != nil {
			return nil, errors.Wrap(err, "task repository")
		}
		if listID == parentListID {
			l := &task.Task{
				ID:           id,
				ListID:       listID,
				ProjectID:    projectID,
				UserID:       userID,
				IssueNumber:  issueNumber,
				Title:        title,
				Description:  description,
				PullRequest:  pullRequest,
				HTMLURL:      htmlURL,
				DisplayIndex: displayIndex,
			}
			result = append(result, l)
		}
	}
	return result, nil
}

// NonIssueTasks returns all tasks related a list.
func (t *Task) NonIssueTasks(projectID, userID int64) ([]*task.Task, error) {
	result := []*task.Task{}
	rows, err := t.db.Query("SELECT id, list_id, project_id, user_id, issue_number, title, description, pull_request, html_url, display_index FROM tasks WHERE project_id = ? and user_id = ? and issue_number IS NULL;", projectID, userID)
	if err != nil {
		return nil, errors.Wrap(err, "task repository")
	}

	for rows.Next() {
		var id, listID, userID, projectID, displayIndex int64
		var title, description string
		var issueNumber sql.NullInt64
		var pullRequest bool
		var htmlURL sql.NullString
		err := rows.Scan(&id, &listID, &projectID, &userID, &issueNumber, &title, &description, &pullRequest, &htmlURL, &displayIndex)
		if err != nil {
			return nil, errors.Wrap(err, "task repository")
		}
		l := &task.Task{
			ID:           id,
			ListID:       listID,
			ProjectID:    projectID,
			UserID:       userID,
			IssueNumber:  issueNumber,
			Title:        title,
			Description:  description,
			PullRequest:  pullRequest,
			HTMLURL:      htmlURL,
			DisplayIndex: displayIndex,
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
func (t *Task) Update(id, listID, projectID, userID int64, issueNumber sql.NullInt64, title, description string, pullRequest bool, htmlURL sql.NullString, displayIndex int64, tx *sql.Tx) error {
	var err error
	if tx != nil {
		_, err = tx.Exec("update tasks set list_id = ?, project_id = ?, user_id = ?, issue_number = ?, title = ?, description = ?, pull_request = ?, html_url = ?, display_index = ? where id = ?;", listID, projectID, userID, issueNumber, title, description, pullRequest, htmlURL, displayIndex, id)
	} else {
		_, err = t.db.Exec("update tasks set list_id = ?, project_id = ?, user_id = ?, issue_number = ?, title = ?, description = ?, pull_request = ?, html_url = ?, display_index = ? where id = ?;", listID, projectID, userID, issueNumber, title, description, pullRequest, htmlURL, displayIndex, id)
	}
	if err != nil {
		return errors.Wrap(err, "task repository")
	}
	logging.SharedInstance().MethodInfo("task", "Update").Debugf("task updated: %+v", t)

	return nil
}

// PushOutAfterTasks updates display_index of after tasks.
func (t *Task) PushOutAfterTasks(listID int64, sinceDisplayIndex int64, tx *sql.Tx) error {
	_, err := tx.Exec("update tasks set display_index = display_index + 1 where id in (select id from (select id from tasks where list_id = ? and display_index >= ?) as tmp);", listID, sinceDisplayIndex)
	return err
}

// GetMaxDisplayIndex gets display index of the last task.
func (t *Task) GetMaxDisplayIndex(listID int64) (*int64, error) {
	// When the list does not have any tasks, max id is null.
	// But null is not error, so we have to accept null value.
	var index interface{}
	err := t.db.QueryRow("select max(display_index) from tasks where list_id = ?;", listID).Scan(&index)
	if err != nil {
		return nil, err
	}
	if index == nil {
		return nil, nil
	}
	displayIndex := index.(int64)
	return &displayIndex, nil
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

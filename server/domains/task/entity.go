package task

import (
	"database/sql"

	"github.com/pkg/errors"
)

// Task has a task model object
type Task struct {
	ID             int64
	ListID         int64
	ProjectID      int64
	UserID         int64
	IssueNumber    sql.NullInt64
	Title          string
	Description    string
	PullRequest    bool
	HTMLURL        sql.NullString
	infrastructure Repository
}

type Repository interface {
	Find(int64) (int64, int64, int64, int64, sql.NullInt64, string, string, bool, sql.NullString, error)
	FindByIssueNumber(int64, int) (int64, int64, int64, int64, sql.NullInt64, string, string, bool, sql.NullString, error)
	Create(int64, int64, int64, sql.NullInt64, string, string, bool, sql.NullString) (int64, error)
	Update(int64, int64, int64, int64, sql.NullInt64, string, string, bool, sql.NullString) error
	ChangeList(int64, int64, *int64) error
	Delete(int64) error
	Tasks(int64) ([]map[string]interface{}, error)
	NonIssueTasks(int64, int64) ([]map[string]interface{}, error)
}

// New returns a task entity
func New(id, listID, projectID, userID int64, issueNumber sql.NullInt64, title, description string, pullRequest bool, htmlURL sql.NullString, infrastructure Repository) *Task {
	return &Task{
		id,
		listID,
		projectID,
		userID,
		issueNumber,
		title,
		description,
		pullRequest,
		htmlURL,
		infrastructure,
	}
}

// Create call save in model
func (t *Task) Create() error {
	id, err := t.infrastructure.Create(t.ListID, t.ProjectID, t.UserID, t.IssueNumber, t.Title, t.Description, t.PullRequest, t.HTMLURL)
	if err != nil {
		return nil
	}
	t.ID = id
	return nil
}

// Update call update in model
func (t *Task) Update(listID int64, issueNumber sql.NullInt64, title, description string, pullRequest bool, htmlURL sql.NullString) error {
	return t.infrastructure.Update(t.ID, listID, t.ProjectID, t.UserID, issueNumber, title, description, pullRequest, htmlURL)
}

// ChangeList change list of a task, and reorder task
// returns isReorder, error.
func (t *Task) ChangeList(listID int64, prevToTaskID *int64) (bool, error) {
	var isReorder bool
	if listID == t.ListID {
		isReorder = true
	} else {
		isReorder = false
	}

	err := t.infrastructure.ChangeList(t.ID, listID, prevToTaskID)
	if err != nil {
		return isReorder, err
	}
	t.ListID = listID
	return isReorder, nil
}

// Delete call delete in model
func (t *Task) Delete() error {
	if t.IssueNumber.Valid {
		return errors.New("task is related issues")
	}
	return t.infrastructure.Delete(t.ID)
}

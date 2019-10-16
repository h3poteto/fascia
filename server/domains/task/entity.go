package task

import (
	"database/sql"
	"errors"
)

// Task has a task model object
type Task struct {
	ID           int64
	ListID       int64
	ProjectID    int64
	UserID       int64
	IssueNumber  sql.NullInt64
	Title        string
	Description  string
	PullRequest  bool
	HTMLURL      sql.NullString
	DisplayIndex int64
}

// New returns a task entity
func New(id, listID, projectID, userID int64, issueNumber sql.NullInt64, title, description string, pullRequest bool, htmlURL sql.NullString) *Task {
	return &Task{
		ID:          id,
		ListID:      listID,
		ProjectID:   projectID,
		UserID:      userID,
		IssueNumber: issueNumber,
		Title:       title,
		Description: description,
		PullRequest: pullRequest,
		HTMLURL:     htmlURL,
	}
}

// Update updates a task.
func (t *Task) Update(listID int64, issueNumber sql.NullInt64, title, description string, pullRequest bool, htmlURL sql.NullString, displayIndex int64) error {
	t.ListID = listID
	t.IssueNumber = issueNumber
	t.Title = title
	t.Description = description
	t.PullRequest = pullRequest
	t.HTMLURL = htmlURL
	t.DisplayIndex = displayIndex
	return nil
}

// Deletable returns whether the task can delete or not.
func (t *Task) Deletable() error {
	if t.IssueNumber.Valid {
		return errors.New("task is related issues")
	}
	return nil
}

package views

import (
	"github.com/h3poteto/fascia/server/domains/entities/task"
)

// Task provides a response structure for task
type Task struct {
	ID          int64  `json:ID`
	ListID      int64  `json:ListID`
	UserID      int64  `json:UserID`
	IssueNumber int64  `json:IssueNumber`
	Title       string `json:Title`
	Description string `json:Description`
	HTMLURL     string `json:HTMLURL`
	PullRequest bool   `json:PullRequest`
}

// ParseTaskJSON returns a Task struct for response
func ParseTaskJSON(task *task.Task) (*Task, error) {
	return &Task{
		ID:          task.ID,
		ListID:      task.ListID,
		UserID:      task.UserID,
		IssueNumber: task.IssueNumber.Int64,
		Title:       task.Title,
		Description: task.Description,
		HTMLURL:     task.HTMLURL.String,
		PullRequest: task.PullRequest,
	}, nil
}

// ParseTasksJSON returns some Task structs for response
func ParseTasksJSON(tasks []*task.Task) ([]*Task, error) {
	results := []*Task{}
	for _, t := range tasks {
		parse, err := ParseTaskJSON(t)
		if err != nil {
			return nil, err
		}
		results = append(results, parse)
	}
	return results, nil
}

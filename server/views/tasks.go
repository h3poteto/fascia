package views

import (
	"github.com/h3poteto/fascia/server/entities/task"
)

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

func ParseTaskJSON(task *task.Task) (*Task, error) {
	return &Task{
		ID:          task.TaskModel.ID,
		ListID:      task.TaskModel.ListID,
		UserID:      task.TaskModel.UserID,
		IssueNumber: task.TaskModel.IssueNumber.Int64,
		Title:       task.TaskModel.Title,
		Description: task.TaskModel.Description,
		HTMLURL:     task.TaskModel.HTMLURL.String,
		PullRequest: task.TaskModel.PullRequest,
	}, nil
}

func ParseTasksJSON(tasks []*task.Task) ([]*Task, error) {
	results := make([]*Task, 0)
	for _, t := range tasks {
		parse, err := ParseTaskJSON(t)
		if err != nil {
			return nil, err
		}
		results = append(results, parse)
	}
	return results, nil
}

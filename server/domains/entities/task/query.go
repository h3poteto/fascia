package task

import (
	"github.com/h3poteto/fascia/server/infrastructures/task"
)

// Find returns a task entity
func Find(listID, taskID int64) (*Task, error) {
	t := &Task{
		ID:     taskID,
		ListID: listID,
	}
	err := t.reload()
	if err != nil {
		return nil, err
	}
	return t, nil
}

// FindByIssueNumber returns a task entity
func FindByIssueNumber(projectID int64, issueNumber int) (*Task, error) {
	infrastructure, err := task.FindByIssueNumber(projectID, issueNumber)
	if err != nil {
		return nil, err
	}
	t := &Task{
		infrastructure: infrastructure,
	}
	err = t.reload()
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Tasks returns all tasks related a list.
func Tasks(listID int64) ([]*Task, error) {
	var slice []*Task

	tasks, err := task.Tasks(listID)
	if err != nil {
		return nil, err
	}
	for _, t := range tasks {
		s := &Task{
			infrastructure: t,
		}
		if err := s.reload(); err != nil {
			return nil, err
		}
		slice = append(slice, s)
	}
	return slice, nil
}

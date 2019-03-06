package task

import (
	"database/sql"
	"errors"
)

// Find returns a task entity
func Find(targetTaskID int64, infrastructure Repository) (*Task, error) {
	id, listID, projectID, userID, issueNumber, title, description, pullRequest, htmlURL, err := infrastructure.Find(targetTaskID)
	if err != nil {
		return nil, err
	}
	return New(id, listID, projectID, userID, issueNumber, title, description, pullRequest, htmlURL, infrastructure), nil
}

// FindByIssueNumber returns a task entity
func FindByIssueNumber(targetProjectID int64, targetIssueNumber int, infrastructure Repository) (*Task, error) {
	id, listID, projectID, userID, issueNumber, title, description, pullRequest, htmlURL, err := infrastructure.FindByIssueNumber(targetProjectID, targetIssueNumber)
	if err != nil {
		return nil, err
	}
	return New(id, listID, projectID, userID, issueNumber, title, description, pullRequest, htmlURL, infrastructure), nil
}

// Tasks returns all tasks related a list.
func Tasks(targetListID int64, infrastructure Repository) ([]*Task, error) {
	var result []*Task

	tasks, err := infrastructure.Tasks(targetListID)
	if err != nil {
		return result, err
	}
	for _, task := range tasks {
		id, ok := task["id"].(int64)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		listID, ok := task["listID"].(int64)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		projectID, ok := task["projectID"].(int64)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		userID, ok := task["userID"].(int64)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		issueNumber, ok := task["issueNumber"].(sql.NullInt64)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		title, ok := task["title"].(string)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		description, ok := task["description"].(string)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		pullRequest, ok := task["pullRequest"].(bool)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		htmlURL, ok := task["htmlURL"].(sql.NullString)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		t := New(id, listID, projectID, userID, issueNumber, title, description, pullRequest, htmlURL, infrastructure)
		result = append(result, t)
	}
	return result, nil
}

// NoIssueTasks returns all tasks related a list.
func NonIssueTasks(targetProjectID, targetUserID int64, infrastructure Repository) ([]*Task, error) {
	var result []*Task

	tasks, err := infrastructure.NonIssueTasks(targetProjectID, targetUserID)
	if err != nil {
		return result, err
	}
	for _, task := range tasks {
		id, ok := task["id"].(int64)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		listID, ok := task["listID"].(int64)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		projectID, ok := task["projectID"].(int64)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		userID, ok := task["userID"].(int64)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		issueNumber, ok := task["issueNumber"].(sql.NullInt64)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		title, ok := task["title"].(string)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		description, ok := task["description"].(string)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		pullRequest, ok := task["pullRequest"].(bool)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		htmlURL, ok := task["htmlURL"].(sql.NullString)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		t := New(id, listID, projectID, userID, issueNumber, title, description, pullRequest, htmlURL, infrastructure)
		result = append(result, t)
	}
	return result, nil
}

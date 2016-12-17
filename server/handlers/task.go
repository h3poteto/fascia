package handlers

import (
	"database/sql"

	"github.com/h3poteto/fascia/server/services"
)

func SaveTask(task *services.Task) error {
	err := task.Save()
	if err != nil {
		return err
	}

	go func(task *services.Task) {
		projectID := task.TaskAggregation.TaskModel.ProjectID
		p, err := services.FindProject(projectID)
		// TODO: log
		if err != nil {
			return
		}
		token, err := p.ProjectAggregation.OauthToken()
		if err != nil {
			return
		}
		repo, err := p.ProjectAggregation.Repository()
		if err != nil {
			return
		}
		err = task.FetchCreated(token, repo)
		if err != nil {
			return
		}
	}(task)

	return nil
}

func UpdateTask(task *services.Task, listID int64, issueNumber sql.NullInt64, title, description string, pullRequest bool, htmlURL sql.NullString) error {
	err := task.Update(listID, issueNumber, title, description, pullRequest, htmlURL)
	if err != nil {
		return err
	}

	go func(task *services.Task) {
		projectID := task.TaskAggregation.TaskModel.ProjectID
		p, err := services.FindProject(projectID)
		// TODO: log
		if err != nil {
			return
		}
		token, err := p.ProjectAggregation.OauthToken()
		if err != nil {
			return
		}
		repo, err := p.ProjectAggregation.Repository()
		if err != nil {
			return
		}
		err = task.FetchUpdated(token, repo)
		if err != nil {
			return
		}
	}(task)
	return nil
}

func ChangeListTask(task *services.Task, listID int64, prevToTaskID *int64) error {
	isReorder, err := task.ChangeList(listID, prevToTaskID)
	if err != nil {
		return err
	}

	go func(task *services.Task, isReorder bool) {
		projectID := task.TaskAggregation.TaskModel.ProjectID
		p, err := services.FindProject(projectID)
		// TODO: log
		if err != nil {
			return
		}
		token, err := p.ProjectAggregation.OauthToken()
		if err != nil {
			return
		}
		repo, err := p.ProjectAggregation.Repository()
		if err != nil {
			return
		}
		err = task.FetchChangedList(token, repo, isReorder)
		if err != nil {
			return
		}
	}(task, isReorder)

	return nil
}

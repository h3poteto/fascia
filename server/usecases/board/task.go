package board

import (
	"database/sql"

	"github.com/h3poteto/fascia/server/domains/services"
	domain "github.com/h3poteto/fascia/server/domains/task"
	repository "github.com/h3poteto/fascia/server/infrastructures/task"
)

// InjectTaskRepository returns a task Repository.
func InjectTaskRepository() domain.Repository {
	return repository.New(InjectDB())
}

// FindTask finds a task.
func FindTask(id int64) (*domain.Task, error) {
	infra := InjectTaskRepository()
	return infra.Find(id)
}

// CreateTask creates a task, and sync to github.
func CreateTask(listID, projectID, userID int64, issueNumber sql.NullInt64, title, description string, pullRequest bool, htmlURL sql.NullString) (*domain.Task, error) {
	task := domain.New(0, listID, projectID, userID, issueNumber, title, description, pullRequest, htmlURL)
	infra := InjectTaskRepository()
	id, err := infra.Create(task.ListID, task.ProjectID, task.UserID, task.IssueNumber, task.Title, task.Description, task.PullRequest, task.HTMLURL)
	if err != nil {
		return nil, err
	}
	task.ID = id

	go services.AfterCreateTask(task, InjectProjectRepository(), InjectListRepository(), InjectTaskRepository(), InjectRepoRepository())

	return task, nil
}

// UpdateTask updates a task, and sync to github.
func UpdateTask(task *domain.Task, listID int64, issueNumber sql.NullInt64, title, description string, pullRequest bool, htmlURL sql.NullString) error {
	err := task.Update(listID, issueNumber, title, description, pullRequest, htmlURL)
	if err != nil {
		return err
	}
	infra := InjectTaskRepository()
	err = infra.Update(
		task.ID,
		task.ListID,
		task.ProjectID,
		task.UserID,
		task.IssueNumber,
		task.Title,
		task.Description,
		task.PullRequest,
		task.HTMLURL,
	)
	if err != nil {
		return err
	}

	go services.AfterUpdateTask(task, InjectProjectRepository(), InjectListRepository(), InjectRepoRepository())

	return err
}

// TaskChangeList changes the task and sync it to github.
func TaskChangeList(task *domain.Task, listID int64, prevToTaskID *int64) error {
	isReorder, err := task.ChangeList(listID, prevToTaskID)
	if err != nil {
		return err
	}
	infra := InjectTaskRepository()
	err = infra.ChangeList(task.ID, listID, prevToTaskID)
	if err != nil {
		return err
	}

	go services.AfterTaskChangeList(task, isReorder, InjectProjectRepository(), InjectListRepository(), InjectRepoRepository())
	return nil
}

// DeleteTask deletes a task.
func DeleteTask(t *domain.Task) error {
	if err := t.Deletable(); err != nil {
		return err
	}
	infra := InjectTaskRepository()
	return infra.Delete(t.ID)
}

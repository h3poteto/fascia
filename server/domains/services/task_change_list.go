package services

import (
	"database/sql"

	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/domains/list"
	"github.com/h3poteto/fascia/server/domains/project"
	"github.com/h3poteto/fascia/server/domains/repo"
	"github.com/h3poteto/fascia/server/domains/task"
)

func TaskInsertMiddle(targetTask *task.Task, listID int64, prevToTaskID int64, taskInfra task.Repository, tx *sql.Tx) (*task.Task, error) {
	prevTask, err := taskInfra.Find(prevToTaskID)
	if err != nil {
		return nil, err
	}
	prevToTaskIndex := prevTask.DisplayIndex
	err = taskInfra.PushOutAfterTasks(listID, prevToTaskIndex, tx)
	if err != nil {
		return nil, err
	}
	targetTask.Update(listID, targetTask.IssueNumber, targetTask.Title, targetTask.Description, targetTask.PullRequest, targetTask.HTMLURL, prevToTaskIndex)
	return targetTask, nil
}

func TaskInsertLast(targetTask *task.Task, listID int64, taskInfra task.Repository) (*task.Task, error) {
	maxIndex, err := taskInfra.GetMaxDisplayIndex(listID)
	if err != nil {
		return nil, err
	}
	displayIndex := int64(1)
	if maxIndex != nil {
		displayIndex = *maxIndex + 1
	}
	targetTask.Update(listID, targetTask.IssueNumber, targetTask.Title, targetTask.Description, targetTask.PullRequest, targetTask.HTMLURL, displayIndex)
	return targetTask, nil
}

// AfterTaskChangeList fetch the changed list.
func AfterTaskChangeList(t *task.Task, isReorder bool, projectInfra project.Repository, listInfra list.Repository, repoInfra repo.Repository) {
	projectID := t.ProjectID
	p, err := projectInfra.Find(projectID)
	// TODO: log
	if err != nil {
		return
	}
	token, err := projectInfra.OauthToken(p.ID)
	if err != nil {
		return
	}
	repo, err := repoInfra.FindByProjectID(p.ID)
	if err != nil {
		return
	}
	err = fetchChangedList(t, token, repo, isReorder, listInfra)
	if err != nil {
		return
	}
}

func fetchChangedList(t *task.Task, oauthToken string, repo *repo.Repo, isReorder bool, listInfra list.Repository) error {
	if !isReorder && repo != nil {
		_, err := syncTaskToIssue(t, repo, oauthToken, listInfra)
		if err != nil {
			logging.SharedInstance().MethodInfo("task", "ChangeList").Error(err)
			return err
		}
		logging.SharedInstance().MethodInfo("Task", "Update").Debugf("task synced to github: %+v", t)
	}
	return nil
}

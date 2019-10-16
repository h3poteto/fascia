package services

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/domains/list"
	"github.com/h3poteto/fascia/server/domains/project"
	"github.com/h3poteto/fascia/server/domains/repo"
	"github.com/h3poteto/fascia/server/domains/task"
)

func AfterUpdateTask(t *task.Task, projectInfra project.Repository, listInfra list.Repository, repoInfra repo.Repository) {
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
	err = fetchUpdatedTask(t, token, repo, listInfra)
	if err != nil {
		return
	}
}

func fetchUpdatedTask(t *task.Task, oauthToken string, repo *repo.Repo, listInfra list.Repository) error {
	// github側へ同期
	if repo != nil {
		_, err := syncTaskToIssue(t, repo, oauthToken, listInfra)
		if err != nil {
			logging.SharedInstance().MethodInfo("task", "Update").Error(err)
			return err
		}
		logging.SharedInstance().MethodInfo("task", "Update").Debugf("task synced to github: %+v", t)
	}
	return nil
}

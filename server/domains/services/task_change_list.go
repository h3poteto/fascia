package services

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/domains/list"
	"github.com/h3poteto/fascia/server/domains/project"
	"github.com/h3poteto/fascia/server/domains/repo"
	"github.com/h3poteto/fascia/server/domains/task"
)

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

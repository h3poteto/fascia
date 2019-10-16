package services

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/domains/list"
	"github.com/h3poteto/fascia/server/domains/project"
	"github.com/h3poteto/fascia/server/domains/repo"
)

// AfterCreateList fetch the created list.
func AfterCreateList(l *list.List, projectInfra project.Repository, repoInfra repo.Repository) {
	projectID := l.ProjectID
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
	err = fetchCreatedList(l, token, repo)
	if err != nil {
		return
	}
}

func fetchCreatedList(l *list.List, oauthToken string, repo *repo.Repo) error {
	if repo != nil {
		label, err := repo.CheckLabelPresent(oauthToken, l.Title.String)
		if err != nil {
			return err
		} else if label == nil {
			label, err = repo.CreateGithubLabel(oauthToken, l.Title.String, l.Color.String)
			if err != nil {
				return err
			}
		} else {
			// 色だけはこちら指定のものに変更したい
			_, err := repo.UpdateGithubLabel(oauthToken, l.Title.String, l.Title.String, l.Color.String)
			if err != nil {
				return err
			}
			logging.SharedInstance().MethodInfo("list", "Save").Info("github label already exist")
		}
	}
	return nil
}

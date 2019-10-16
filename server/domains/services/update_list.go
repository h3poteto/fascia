package services

import (
	"github.com/h3poteto/fascia/server/domains/list"
	"github.com/h3poteto/fascia/server/domains/project"
	"github.com/h3poteto/fascia/server/domains/repo"
)

func AfterUpdateList(l *list.List, title, color string, projectInfra project.Repository, repoInfra repo.Repository) {
	projectID := l.ProjectID
	p, err := projectInfra.Find(projectID)
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
	err = fetchUpdatedList(l, token, repo, title, color)
	if err != nil {
		return
	}
}

func fetchUpdatedList(l *list.List, oauthToken string, repo *repo.Repo, newTitle, newColor string) error {
	if repo != nil {
		// 編集前のラベルがそもそも存在しているかどうかを確認する
		existLabel, err := repo.CheckLabelPresent(oauthToken, l.Title.String)
		if err != nil {
			return err
		} else if existLabel == nil {
			// editの場合ここに入る可能性は，createのgithub同期がうまく動いていない場合のみ
			// 編集前のラベルが存在しなければ新しく作る
			_, err := repo.CreateGithubLabel(oauthToken, newTitle, newColor)
			if err != nil {
				return err
			}
		} else {
			_, err := repo.UpdateGithubLabel(oauthToken, l.Title.String, newTitle, newColor)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

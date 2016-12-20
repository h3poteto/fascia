package services

import (
	"database/sql"

	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/aggregations/list"
	"github.com/h3poteto/fascia/server/aggregations/repository"
)

type List struct {
	ListAggregation *list.List
}

func NewList(id, projectID, userID int64, title, color string, optionID sql.NullInt64, isHidden bool) *List {
	return &List{
		ListAggregation: list.New(id, projectID, userID, title, color, optionID, isHidden),
	}
}

func FindListByID(projectID, listID int64) (*List, error) {
	l, err := list.FindByID(projectID, listID)
	if err != nil {
		return nil, err
	}
	return &List{
		ListAggregation: l,
	}, nil
}

func (l *List) Save() error {
	err := l.ListAggregation.Save(nil)
	if err != nil {
		return err
	}
	go func(list *List) {
		projectID := list.ListAggregation.ListModel.ProjectID
		p, err := FindProject(projectID)
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
		err = list.fetchCreated(token, repo)
		if err != nil {
			return
		}
	}(l)
	return nil
}

func (l *List) fetchCreated(oauthToken string, repo *repository.Repository) error {
	if repo != nil {
		label, err := repo.CheckLabelPresent(oauthToken, l.ListAggregation.ListModel.Title.String)
		if err != nil {
			return err
		} else if label == nil {
			label, err = repo.CreateGithubLabel(oauthToken, l.ListAggregation.ListModel.Title.String, l.ListAggregation.ListModel.Color.String)
			if err != nil {
				return err
			}
		} else {
			// 色だけはこちら指定のものに変更したい
			_, err := repo.UpdateGithubLabel(oauthToken, l.ListAggregation.ListModel.Title.String, l.ListAggregation.ListModel.Title.String, l.ListAggregation.ListModel.Color.String)
			if err != nil {
				return err
			}
			logging.SharedInstance().MethodInfo("list", "Save").Info("github label already exist")
		}
	}
	return nil
}

func (l *List) Update(title, color string, optionID int64) error {
	err := l.ListAggregation.UpdateExceptInitList(title, color, optionID)
	if err != nil {
		return err
	}

	go func(list *List, title, color string) {
		projectID := list.ListAggregation.ListModel.ProjectID
		p, err := FindProject(projectID)
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
		err = list.fetchUpdated(token, repo, title, color)
		if err != nil {
			return
		}
	}(l, title, color)
	return nil
}

func (l *List) fetchUpdated(oauthToken string, repo *repository.Repository, newTitle, newColor string) error {
	if repo != nil {
		// 編集前のラベルがそもそも存在しているかどうかを確認する
		existLabel, err := repo.CheckLabelPresent(oauthToken, l.ListAggregation.ListModel.Title.String)
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
			_, err := repo.UpdateGithubLabel(oauthToken, l.ListAggregation.ListModel.Title.String, newTitle, newColor)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *List) Hide() error {
	return l.ListAggregation.Hide()
}

func (l *List) Display() error {
	return l.ListAggregation.Display()
}

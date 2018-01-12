package services

import (
	"database/sql"

	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/entities/list"
	"github.com/h3poteto/fascia/server/entities/repository"
)

// List has a list entity
type List struct {
	ListEntity *list.List
}

// NewList returns a list service
func NewList(id, projectID, userID int64, title, color string, optionID sql.NullInt64, isHidden bool) *List {
	return &List{
		ListEntity: list.New(id, projectID, userID, title, color, optionID, isHidden),
	}
}

// FindListByID returns a list service
func FindListByID(projectID, listID int64) (*List, error) {
	l, err := list.FindByID(projectID, listID)
	if err != nil {
		return nil, err
	}
	return &List{
		ListEntity: l,
	}, nil
}

// Save save list entity, and fetch created list to github
func (l *List) Save() error {
	err := l.ListEntity.Save(nil)
	if err != nil {
		return err
	}
	go func(list *List) {
		projectID := list.ListEntity.ProjectID
		p, err := FindProject(projectID)
		// TODO: log
		if err != nil {
			return
		}
		token, err := p.ProjectEntity.OauthToken()
		if err != nil {
			return
		}
		repo, err := p.ProjectEntity.Repository()
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
		label, err := repo.CheckLabelPresent(oauthToken, l.ListEntity.Title.String)
		if err != nil {
			return err
		} else if label == nil {
			label, err = repo.CreateGithubLabel(oauthToken, l.ListEntity.Title.String, l.ListEntity.Color.String)
			if err != nil {
				return err
			}
		} else {
			// 色だけはこちら指定のものに変更したい
			_, err := repo.UpdateGithubLabel(oauthToken, l.ListEntity.Title.String, l.ListEntity.Title.String, l.ListEntity.Color.String)
			if err != nil {
				return err
			}
			logging.SharedInstance().MethodInfo("list", "Save").Info("github label already exist")
		}
	}
	return nil
}

// Update save list entity, and fetch updated list to github
func (l *List) Update(title, color string, optionID int64) error {
	err := l.ListEntity.UpdateExceptInitList(title, color, optionID)
	if err != nil {
		return err
	}

	go func(list *List, title, color string) {
		projectID := list.ListEntity.ProjectID
		p, err := FindProject(projectID)
		if err != nil {
			return
		}
		token, err := p.ProjectEntity.OauthToken()
		if err != nil {
			return
		}
		repo, err := p.ProjectEntity.Repository()
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
		existLabel, err := repo.CheckLabelPresent(oauthToken, l.ListEntity.Title.String)
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
			_, err := repo.UpdateGithubLabel(oauthToken, l.ListEntity.Title.String, newTitle, newColor)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Hide change list visibility
func (l *List) Hide() error {
	return l.ListEntity.Hide()
}

// Display change list visibility
func (l *List) Display() error {
	return l.ListEntity.Display()
}

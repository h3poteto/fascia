package services

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/aggregations/list"
	"github.com/h3poteto/fascia/server/aggregations/repository"
)

type List struct {
	ListAggregation *list.List
}

func (l *List) Save() error {
	return l.ListAggregation.Save()
}

func (l *List) FetchCreated(oauthToken string, repo *repository.Repository) error {
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
	return l.ListAggregation.Update(title, color, optionID)
}

func (l *List) FetchUpdated(oauthToken string, repo *repository.Repository, newTitle, newColor string) error {
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

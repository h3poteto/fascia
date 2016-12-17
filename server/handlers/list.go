package handlers

import (
	"github.com/h3poteto/fascia/server/services"
)

func SaveList(list *services.List) error {
	err := list.Save()
	if err != nil {
		return err
	}

	go func(list *services.List) {
		projectID := list.ListAggregation.ListModel.ProjectID
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
		err = list.FetchCreated(token, repo)
		if err != nil {
			return
		}
	}(list)

	return nil
}

func UpdateList(list *services.List, title, color string, listOptionID int64) error {
	err := list.Update(title, color, listOptionID)
	if err != nil {
		return err
	}

	go func(list *services.List, title, color string) {
		projectID := list.ListAggregation.ListModel.ProjectID
		p, err := services.FindProject(projectID)
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
		err = list.FetchUpdated(token, repo, title, color)
		if err != nil {
			return
		}
	}(list, title, color)
	return nil
}

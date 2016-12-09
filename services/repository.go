package services

import (
	"github.com/h3poteto/fascia/models/repository"
	"github.com/h3poteto/fascia/modules/hub"
)

// CreateRepository create repository record based on github repository
func CreateRepository(ID int, oauthToken string) (*repository.RepositoryStruct, error) {
	// confirm github
	h := hub.New(oauthToken)
	githubRepo, err := h.GetRepository(ID)
	if err != nil {
		return nil, err
	}
	// generate webhook key
	key := repository.GenerateWebhookKey(*githubRepo.Name)
	// save
	repo := repository.New(0, int64(*githubRepo.ID), *githubRepo.Owner.Login, *githubRepo.Name, key)
	err = repo.Save()
	if err != nil {
		return nil, err
	}
	return repo, nil
}

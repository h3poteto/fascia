package services

import (
	"github.com/h3poteto/fascia/server/entities/repository"
)

// Repository has a repository entity
type Repository struct {
	RepositoryEntity *repository.Repository
}

// FindRepositoryByGithubRepoID search repository according to github repository id
func FindRepositoryByGithubRepoID(id int64) (*Repository, error) {
	r, err := repository.FindByGithubRepoID(id)
	if err != nil {
		return nil, err
	}
	return &Repository{
		RepositoryEntity: r,
	}, nil
}

// Authenticate is check token and webhookKey with response
func (r *Repository) Authenticate(token string, response []byte) error {
	return r.RepositoryEntity.Authenticate(token, response)
}

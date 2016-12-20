package services

import (
	"github.com/h3poteto/fascia/server/entities/repository"
)

type Repository struct {
	RepositoryEntity *repository.Repository
}

func FindRepositoryByGithubRepoID(id int64) (*Repository, error) {
	r, err := repository.FindByGithubRepoID(id)
	if err != nil {
		return nil, err
	}
	return &Repository{
		RepositoryEntity: r,
	}, nil
}

func (r *Repository) Authenticate(token string, response []byte) error {
	return r.RepositoryEntity.Authenticate(token, response)
}

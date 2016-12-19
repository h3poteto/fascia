package services

import (
	"github.com/h3poteto/fascia/server/aggregations/repository"
)

type Repository struct {
	RepositoryAggregation *repository.Repository
}

func FindRepositoryByGithubRepoID(id int64) (*Repository, error) {
	r, err := repository.FindByGithubRepoID(id)
	if err != nil {
		return nil, err
	}
	return &Repository{
		RepositoryAggregation: r,
	}, nil
}

func (r *Repository) Authenticate(token string, response []byte) error {
	return r.RepositoryAggregation.Authenticate(token, response)
}

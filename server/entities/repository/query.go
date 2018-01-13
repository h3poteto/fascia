package repository

import (
	"github.com/h3poteto/fascia/server/infrastructures/repository"
)

// FindByGithubRepoID find repository entity according to repository id in github
func FindByGithubRepoID(id int64) (*Repository, error) {
	r := &Repository{
		RepositoryID: id,
	}
	if err := r.reload(); err != nil {
		return nil, err
	}
	return r, nil
}

// FindByProjectID returns a repository related a project.
func FindByProjectID(projectID int64) (*Repository, error) {
	infrastructure, err := repository.FindByProjectID(projectID)
	if err != nil {
		return nil, err
	}
	r := &Repository{
		infrastructure: infrastructure,
	}
	if err := r.reload(); err != nil {
		return nil, err
	}
	return r, nil
}

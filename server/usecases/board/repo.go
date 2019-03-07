package board

import (
	repo "github.com/h3poteto/fascia/server/domains/repo"
	repository "github.com/h3poteto/fascia/server/infrastructures/repo"
)

// InjectRepoRepository returns a repo Repository.
func InjectRepoRepository() repo.Repository {
	return repository.New(InjectDB())
}

// FindRepositoryByGithubRepoID search repository according to github repository id
func FindRepositoryByGithubRepoID(id int64) (*repo.Repo, error) {
	return repo.FindByGithubRepoID(id, InjectRepoRepository())
}

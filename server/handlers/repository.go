package handlers

import (
	"github.com/h3poteto/fascia/server/services"
)

// FindRepositoryByGithubRepoID search repository according to github repository id
func FindRepositoryByGithubRepoID(id int64) (*services.Repository, error) {
	return services.FindRepositoryByGithubRepoID(id)
}

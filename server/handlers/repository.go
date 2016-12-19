package handlers

import (
	"github.com/h3poteto/fascia/server/services"
)

func FindRepositoryByGithubRepoID(id int64) (*services.Repository, error) {
	return services.FindRepositoryByGithubRepoID(id)
}

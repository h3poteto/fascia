package handlers

import (
	"github.com/h3poteto/fascia/server/commands/project"
)

// FindRepositoryByGithubRepoID search repository according to github repository id
func FindRepositoryByGithubRepoID(id int64) (*project.Repository, error) {
	return project.FindRepositoryByGithubRepoID(id)
}

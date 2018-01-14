package handlers

import (
	"github.com/h3poteto/fascia/server/commands/board"
)

// FindRepositoryByGithubRepoID search repository according to github repository id
func FindRepositoryByGithubRepoID(id int64) (*board.Repository, error) {
	return board.FindRepositoryByGithubRepoID(id)
}

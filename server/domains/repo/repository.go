package repo

import "database/sql"

// Repository defines repository interface.
type Repository interface {
	FindByGithubRepoID(int64) (*Repo, error)
	FindByProjectID(int64) (*Repo, error)
	Create(int64, sql.NullString, sql.NullString, string) (int64, error)
}

// NotFoundError is an error when repository not found in DB
type NotFoundError struct {
	Err error
}

// Error interface method for error
func (n *NotFoundError) Error() string {
	return n.Err.Error()
}

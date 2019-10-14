package repo

import "database/sql"

// Repository defines repository interface.
type Repository interface {
	FindByGithubRepoID(int64) (*Repo, error)
	FindByProjectID(int64) (*Repo, error)
	Create(int64, sql.NullString, sql.NullString, string) (int64, error)
}

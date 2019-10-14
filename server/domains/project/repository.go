package project

import "database/sql"

// Repository defines repository interface.
type Repository interface {
	Find(int64) (*Project, error)
	FindByRepositoryID(int64) ([]*Project, error)
	Create(int64, string, string, sql.NullInt64, bool, bool, *sql.Tx) (int64, error)
	Update(int64, int64, string, string, sql.NullInt64, bool, bool) error
	Delete(int64) error
	Projects(int64) ([]*Project, error)
	OauthToken(int64) (string, error)
}

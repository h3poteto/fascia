package user

import "database/sql"

// Repository defines repository interface.
type Repository interface {
	Find(int64) (*User, error)
	FindByEmail(string) (*User, error)
	Create(string, string, sql.NullString, sql.NullString, sql.NullInt64, sql.NullString, sql.NullString) (int64, error)
	Update(int64, string, sql.NullString, sql.NullString, sql.NullInt64, sql.NullString, sql.NullString) error
	UpdatePassword(int64, string) error
}

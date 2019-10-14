package task

import "database/sql"

// Repository defines repository interface.
type Repository interface {
	Find(int64) (*Task, error)
	FindByIssueNumber(int64, int) (*Task, error)
	Create(int64, int64, int64, sql.NullInt64, string, string, bool, sql.NullString) (int64, error)
	Update(int64, int64, int64, int64, sql.NullInt64, string, string, bool, sql.NullString) error
	ChangeList(int64, int64, *int64) error
	Delete(int64) error
	Tasks(int64) ([]*Task, error)
	NonIssueTasks(int64, int64) ([]*Task, error)
}

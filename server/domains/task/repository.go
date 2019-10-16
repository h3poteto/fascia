package task

import "database/sql"

// Repository defines repository interface.
type Repository interface {
	Find(int64) (*Task, error)
	FindByIssueNumber(int64, int) (*Task, error)
	Create(int64, int64, int64, sql.NullInt64, string, string, bool, sql.NullString) (int64, error)
	Update(int64, int64, int64, int64, sql.NullInt64, string, string, bool, sql.NullString, int64, *sql.Tx) error
	PushOutAfterTasks(int64, int64, *sql.Tx) error
	Delete(int64) error
	Tasks(int64) ([]*Task, error)
	NonIssueTasks(int64, int64) ([]*Task, error)
	GetMaxDisplayIndex(int64) (*int64, error)
}

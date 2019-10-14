package list

import "database/sql"

// Repository defines repository interface.
type Repository interface {
	Find(int64, int64) (*List, error)
	FindByTaskID(int64) (*List, error)
	Lists(int64) ([]*List, error)
	NoneList(int64) (*List, error)
	FindOptionByAction(string) (*Option, error)
	FindOptionByID(int64) (*Option, error)
	AllOption() ([]*Option, error)
	Create(int64, int64, sql.NullString, sql.NullString, sql.NullInt64, bool, *sql.Tx) (int64, error)
	Update(*List) error
	Delete(int64) error
	DeleteTasks(int64) error
}

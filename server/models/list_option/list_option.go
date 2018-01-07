package list_option

import (
	"database/sql"

	"github.com/h3poteto/fascia/lib/modules/database"

	"github.com/pkg/errors"
)

// ListOption has list option record
type ListOption struct {
	ID     int64
	Action string
	db     *sql.DB
}

// New returns a new list option object
func New(id int64, action string) *ListOption {
	listOption := &ListOption{ID: id, Action: action}
	listOption.initialize()
	return listOption
}

// FindByAction search a list option according to action
func FindByAction(action string) (*ListOption, error) {
	db := database.SharedInstance().Connection

	var listOptionID int64
	err := db.QueryRow("select id from list_options where action = ?;", action).Scan(&listOptionID)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	return New(listOptionID, action), nil
}

// FindByID search a list option according to id
func FindByID(id sql.NullInt64) (*ListOption, error) {
	db := database.SharedInstance().Connection

	if !id.Valid {
		return nil, errors.New("id is not valid")
	}
	var action string
	err := db.QueryRow("select action from list_options where id = ?;", id).Scan(&action)

	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	return New(id.Int64, action), nil
}

func (l *ListOption) initialize() {
	l.db = database.SharedInstance().Connection
}

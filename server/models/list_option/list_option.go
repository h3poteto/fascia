package list_option

import (
	"database/sql"

	"github.com/h3poteto/fascia/server/models/db"

	"github.com/pkg/errors"
)

type ListOption struct {
	ID       int64
	Action   string
	database *sql.DB
}

func New(id int64, action string) *ListOption {
	listOption := &ListOption{ID: id, Action: action}
	listOption.initialize()
	return listOption
}

func FindByAction(action string) (*ListOption, error) {
	database := db.SharedInstance().Connection

	var listOptionID int64
	err := database.QueryRow("select id from list_options where action = ?;", action).Scan(&listOptionID)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	return New(listOptionID, action), nil
}

func FindByID(id sql.NullInt64) (*ListOption, error) {
	database := db.SharedInstance().Connection

	if !id.Valid {
		return nil, errors.New("id is not valid")
	}
	var action string
	err := database.QueryRow("select action from list_options where id = ?;", id).Scan(&action)

	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	return New(id.Int64, action), nil
}

func (u *ListOption) initialize() {
	u.database = db.SharedInstance().Connection
}

package list_option

import (
	"database/sql"

	"github.com/h3poteto/fascia/server/models/db"

	"github.com/pkg/errors"
)

type ListOiption interface {
}

type ListOptionStruct struct {
	ID       int64
	Action   string
	database *sql.DB
}

func NewListOption(id int64, action string) *ListOptionStruct {
	listOption := &ListOptionStruct{ID: id, Action: action}
	listOption.Initialize()
	return listOption
}

// ListOptionAll list up all options
func ListOptionAll() ([]*ListOptionStruct, error) {
	database := db.SharedInstance().Connection
	var slice []*ListOptionStruct
	var id int64
	var action string
	rows, err := database.Query("select id, action from list_options;")
	if err != nil {
		return slice, errors.Wrap(err, "sql select error")
	}
	for rows.Next() {
		err = rows.Scan(&id, &action)
		if err != nil {
			return nil, errors.Wrap(err, "sql select error")
		}
		l := NewListOption(id, action)
		slice = append(slice, l)
	}
	return slice, nil
}

func FindByAction(action string) (*ListOptionStruct, error) {
	database := db.SharedInstance().Connection

	var listOptionID int64
	err := database.QueryRow("select id from list_options where action = ?;", action).Scan(&listOptionID)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	return NewListOption(listOptionID, action), nil
}

func FindByID(id sql.NullInt64) (*ListOptionStruct, error) {
	database := db.SharedInstance().Connection

	if !id.Valid {
		return nil, errors.New("id is not valid")
	}
	var action string
	err := database.QueryRow("select action from list_options where id = ?;", id).Scan(&action)

	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	return NewListOption(id.Int64, action), nil
}

func (u *ListOptionStruct) Initialize() {
	u.database = db.SharedInstance().Connection
}

// CloseAction return whether it is close option
func (u *ListOptionStruct) CloseAction() bool {
	if u.Action == "close" {
		return true
	}
	return false
}

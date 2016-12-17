package list_option

import (
	"database/sql"

	"github.com/h3poteto/fascia/server/models/db"
	"github.com/h3poteto/fascia/server/models/list_option"

	"github.com/pkg/errors"
)

type ListOption struct {
	ListOptionModel *list_option.ListOption
	database        *sql.DB
}

func New(id int64, action string) *ListOption {
	return &ListOption{
		ListOptionModel: list_option.New(id, action),
		database:        db.SharedInstance().Connection,
	}
}

// ListOptionAll list up all options
func ListOptionAll() ([]*ListOption, error) {
	database := db.SharedInstance().Connection
	var slice []*ListOption
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
		l := New(id, action)
		slice = append(slice, l)
	}
	return slice, nil
}

func FindByID(id int64) (*ListOption, error) {
	option, err := list_option.FindByID(sql.NullInt64{Int64: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return &ListOption{
		ListOptionModel: option,
		database:        db.SharedInstance().Connection,
	}, nil
}

func FindByAction(action string) (*ListOption, error) {
	option, err := list_option.FindByAction(action)
	if err != nil {
		return nil, err
	}
	return &ListOption{
		ListOptionModel: option,
		database:        db.SharedInstance().Connection,
	}, nil
}

// IsCloseAction return whether it is close option
func (l *ListOption) IsCloseAction() bool {
	if l.ListOptionModel.Action == "close" {
		return true
	}
	return false
}

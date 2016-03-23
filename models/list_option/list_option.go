package list_option

import (
	"../../modules/logging"
	"../db"
	"database/sql"
	"errors"
)

type ListOiption interface {
}

type ListOptionStruct struct {
	ID       int64
	Action   string
	database db.DB
}

func NewListOption(id int64, action string) *ListOptionStruct {
	listOption := &ListOptionStruct{ID: id, Action: action}
	listOption.Initialize()
	return listOption
}

func ListOptionAll() []*ListOptionStruct {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	var slice []*ListOptionStruct
	var id int64
	var action string
	rows, err := table.Query("select id, action from list_options;")
	if err != nil {
		logging.SharedInstance().MethodInfo("ListOption", "ListOptionAll", true).Panic(err)
		return slice
	}
	for rows.Next() {
		err = rows.Scan(&id, &action)
		if err != nil {
			panic(err)
		}
		l := NewListOption(id, action)
		slice = append(slice, l)
	}
	return slice
}

func FindByAction(action string) (*ListOptionStruct, error) {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	var listOptionID int64
	err := table.QueryRow("select id from list_options where action = ?;", action).Scan(&listOptionID)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListOption", "FindByAction").Info("cannot find list_option")
		return nil, err
	}
	return NewListOption(listOptionID, action), nil
}

func FindByID(id sql.NullInt64) (*ListOptionStruct, error) {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	if id.Valid {
		var action string
		err := table.QueryRow("select action from list_options where id = ?;", id).Scan(&action)

		if err != nil {
			logging.SharedInstance().MethodInfo("ListOption", "FindByID").Infof("cannot find list_option: %v", id)
			return nil, err
		}
		return NewListOption(id.Int64, action), nil
	} else {
		return nil, errors.New("id is not valid")
	}
}

func (u *ListOptionStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}

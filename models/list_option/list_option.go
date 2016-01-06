package list_option

import (
	"../../modules/logging"
	"../db"
)

type ListOiption interface {
}

type ListOptionStruct struct {
	Id       int64
	Action   string
	database db.DB
}

func NewListOption(id int64, action string) *ListOptionStruct {
	listOption := &ListOptionStruct{Id: id, Action: action}
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
	rows, _ := table.Query("select id, action from list_options;")
	for rows.Next() {
		err := rows.Scan(&id, &action)
		if err != nil {
			panic(err.Error())
		}
		l := NewListOption(id, action)
		slice = append(slice, l)
	}
	return slice
}

func FindByAction(action string) *ListOptionStruct {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	var listOptionId int64
	err := table.QueryRow("select id from list_options where action = ?;", action).Scan(&listOptionId)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListOption", "FindByAction").Info("cannot find list_option")
		return nil
	}
	return NewListOption(listOptionId, action)
}

func (u *ListOptionStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}

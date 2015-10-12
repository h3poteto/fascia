package task

import (
	"../db"
	"database/sql"
)

type Task interface {
	Save() bool
}

type TaskStruct struct {
	Id int64
	ListId int64
	Title sql.NullString
	database db.DB
}

func NewTask(id int64, listID int64, title string) *TaskStruct {
	if listID == 0 {
		return nil
	}
	nullTitle := sql.NullString{String: title, Valid: true}
	task := &TaskStruct{Id: id, ListId: listID, Title: nullTitle}
	task.Initialize()
	return task
}

func (u *TaskStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}


func (u *TaskStruct) Save() bool {
	table := u.database.Init()
	defer table.Close()

	result, err := table.Exec("insert into tasks (list_id, title, created_at) values (?, ?, now());", u.ListId, u.Title)
	if err != nil {
		return false
	}
	u.Id, _ = result.LastInsertId()
	return true
}

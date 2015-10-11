package list

import (
	"../db"
	"database/sql"
)

type List interface {
	Save() bool
}

type ListStruct struct {
	Id int64
	ProjectId int64
	Title sql.NullString
	database db.DB
}

func NewList(id int64, projectID int64, title string) *ListStruct {
	if projectID == 0 {
		return nil
	}
	nullTitle := sql.NullString{String: title, Valid: true}
	list := &ListStruct{Id: id, ProjectId: projectID, Title: nullTitle}
	list.Initialize()
	return list
}

func (u *ListStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}

func (u *ListStruct) Save() bool {
	table := u.database.Init()
	defer table.Close()

	result, err := table.Exec("insert into lists (project_id, title, created_at) values (?, ?, now());", u.ProjectId, u.Title)
	if err != nil {
		return false
	}
	u.Id, _ = result.LastInsertId()
	return true
}

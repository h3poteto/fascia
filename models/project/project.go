package project

import (
	"../db"
)

type Project interface {
	Save() bool
}

type ProjectStruct struct {
	Id int64
	UserId int64
	Title string
	database db.DB
}

func NewProject(id int64, user_id int64, title string) *ProjectStruct {
	project := &ProjectStruct{Id: id, UserId: user_id, Title: title}
	project.Initialize()
	return project
}

func (u *ProjectStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}

func (u *ProjectStruct) Save() bool {
	table := u.database.Init()
	defer table.Close()

	result, err := table.Exec("insert into projects (user_id, title, created_at) values (?, ?, now());", u.UserId, u.Title)
	if err != nil {
		return false
	}
	u.Id, _ = result.LastInsertId()
	return true
}

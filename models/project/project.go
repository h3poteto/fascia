package project

import (
	"../db"
	"database/sql"
)

type Project interface {
	Save() bool
}

type ProjectStruct struct {
	Id int64
	UserId sql.NullInt64
	Title string
	database db.DB
}

func NewProject(id int64, userID int64, title string) *ProjectStruct {
	if userID == 0 {
		return nil
	}
	nullUserID := sql.NullInt64{Int64: int64(userID), Valid: true}
	project := &ProjectStruct{Id: id, UserId: nullUserID, Title: title}
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

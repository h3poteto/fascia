package project

import (
	"../db"
	"database/sql"
	"../list"
)

type Project interface {
	Lists() []*list.ListStruct
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

func FindProject(projectID int64) *ProjectStruct {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	var id int64
	var userID sql.NullInt64
	var title string
	rows, _ := table.Query("select id, user_id, title from projects where id = ?;", projectID)
	for rows.Next() {
		err := rows.Scan(&id, &userID, &title)
		if err != nil {
			panic(err.Error())
		}
	}
	if userID.Valid {
		project := NewProject(id, userID.Int64, title)
		return project
	} else {
		return nil
	}
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

func (u *ProjectStruct) Lists() []*list.ListStruct {
	table := u.database.Init()
	defer table.Close()

	rows, _ := table.Query("select id, project_id, title from lists where project_id = ?;", u.Id)
	var slice []*list.ListStruct
	for rows.Next() {
		var id, projectID int64
		var title sql.NullString
		err := rows.Scan(&id, &projectID, &title)
		if err != nil {
			panic(err.Error())
		}
		if projectID == u.Id && title.Valid {
			l := list.NewList(id, projectID, title.String)
			slice = append(slice, l)
		}
	}
	return slice
}

package project

import (
	"../db"
	"database/sql"
	"../list"
	"../repository"
)

type Project interface {
	Lists() []*list.ListStruct
	Save() bool
}

type ProjectStruct struct {
	Id int64
	UserId sql.NullInt64
	Title string
	Description string
	database db.DB
}

func NewProject(id int64, userID int64, title string, description string) *ProjectStruct {
	if userID == 0 {
		return nil
	}
	nullUserID := sql.NullInt64{Int64: int64(userID), Valid: true}
	project := &ProjectStruct{Id: id, UserId: nullUserID, Title: title, Description: description}
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
	var description string
	rows, _ := table.Query("select id, user_id, title, description from projects where id = ?;", projectID)
	for rows.Next() {
		err := rows.Scan(&id, &userID, &title, &description)
		if err != nil {
			panic(err.Error())
		}
	}
	if userID.Valid {
		project := NewProject(id, userID.Int64, title, description)
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

	result, err := table.Exec("insert into projects (user_id, title, description, created_at) values (?, ?, ?, now());", u.UserId, u.Title, u.Description)
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

func (u *ProjectStruct) Repository() *repository.RepositoryStruct {
	table := u.database.Init()
	defer table.Close()

	rows, _ := table.Query("select id, project_id, repository_id, name, owner from repositories where project_id = ?", u.Id)
	for rows.Next() {
		var id, projectId, repositoryId int64
		var owner, name sql.NullString
		err := rows.Scan(&id, &projectId, &repositoryId, &owner, &name)
		if err != nil {
			panic(err.Error())
		}
		if projectId == u.Id && owner.Valid {
			r := repository.NewRepository(id, projectId, repositoryId, owner.String, name.String)
			return r
		}
	}
	return nil
}

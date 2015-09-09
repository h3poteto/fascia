package project

import (
	"../db"
)

type Project interface {
}

type ProjectStruct struct {
	Id int
	Title string
	database db.DB
}

func NewProject(title string) *ProjectStruct {
	project := &ProjectStruct{Title: title}
	project.Initialize()
	return project
}

func (u *ProjectStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}

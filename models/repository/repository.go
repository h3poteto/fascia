package repository

import (
	"../../modules/logging"
	"../db"
	"database/sql"
)

type Repository interface {
	Save() bool
}

type RepositoryStruct struct {
	Id           int64
	RepositoryId int64
	Owner        sql.NullString
	Name         sql.NullString
	database     db.DB
}

func NewRepository(id int64, repositoryId int64, owner string, name string) *RepositoryStruct {
	if repositoryId <= 0 {
		return nil
	}
	repository := &RepositoryStruct{Id: id, RepositoryId: repositoryId, Owner: sql.NullString{String: owner, Valid: true}, Name: sql.NullString{String: name, Valid: true}}
	repository.Initialize()
	return repository
}

func (u *RepositoryStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}

func (u *RepositoryStruct) Save() bool {
	table := u.database.Init()
	defer table.Close()

	result, err := table.Exec("insert into repositories (repository_id, owner, name, created_at) values (?, ?, ?, now());", u.RepositoryId, u.Owner, u.Name)
	if err != nil {
		logging.SharedInstance().MethodInfo("Repository", "Save", true).Errorf("repository save failed: %v", err)
		return false
	}
	u.Id, _ = result.LastInsertId()
	return true
}

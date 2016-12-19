package repository

import (
	"database/sql"

	"github.com/h3poteto/fascia/server/models/db"

	"github.com/pkg/errors"
)

type Repository interface {
	Save() bool
}

type RepositoryStruct struct {
	ID           int64
	RepositoryID int64
	Owner        sql.NullString
	Name         sql.NullString
	WebhookKey   string
	database     *sql.DB
}

// New is build new Repository struct
func New(id int64, repositoryID int64, owner string, name string, webhookKey string) *RepositoryStruct {
	if repositoryID <= 0 {
		return nil
	}
	repository := &RepositoryStruct{ID: id, RepositoryID: repositoryID, Owner: sql.NullString{String: owner, Valid: true}, Name: sql.NullString{String: name, Valid: true}, WebhookKey: webhookKey}
	repository.Initialize()
	return repository
}

// FindByGithubRepoID is return a Repository struct from repository_id
func FindByGithubRepoID(githubRepoID int64) (*RepositoryStruct, error) {
	database := db.SharedInstance().Connection
	var id, repoID int64
	var owner, name, webhookKey string
	err := database.QueryRow("select id, repository_id, owner, name, webhook_key from repositories where repository_id = ?;", githubRepoID).Scan(&id, &repoID, &owner, &name, &webhookKey)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	return New(id, repoID, owner, name, webhookKey), nil
}

func (u *RepositoryStruct) Initialize() {
	u.database = db.SharedInstance().Connection
}

func (u *RepositoryStruct) Save() error {
	result, err := u.database.Exec("insert into repositories (repository_id, owner, name, webhook_key, created_at) values (?, ?, ?, ?, now());", u.RepositoryID, u.Owner, u.Name, u.WebhookKey)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	u.ID, _ = result.LastInsertId()
	return nil
}

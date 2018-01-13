package repository

import (
	"database/sql"

	"github.com/h3poteto/fascia/lib/modules/database"

	"github.com/pkg/errors"
)

// Repository has repository record
type Repository struct {
	ID           int64
	RepositoryID int64
	Owner        sql.NullString
	Name         sql.NullString
	WebhookKey   string
	db           *sql.DB
}

// New is build new Repository struct
func New(id int64, repositoryID int64, owner string, name string, webhookKey string) *Repository {
	if repositoryID <= 0 {
		return nil
	}
	repository := &Repository{ID: id, RepositoryID: repositoryID, Owner: sql.NullString{String: owner, Valid: true}, Name: sql.NullString{String: name, Valid: true}, WebhookKey: webhookKey}
	repository.initialize()
	return repository
}

// FindByGithubRepoID is return a Repository struct from repository_id
func FindByGithubRepoID(githubRepoID int64) (*Repository, error) {
	db := database.SharedInstance().Connection
	var id, repoID int64
	var owner, name, webhookKey string
	err := db.QueryRow("select id, repository_id, owner, name, webhook_key from repositories where repository_id = ?;", githubRepoID).Scan(&id, &repoID, &owner, &name, &webhookKey)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	return New(id, repoID, owner, name, webhookKey), nil
}

func (r *Repository) initialize() {
	r.db = database.SharedInstance().Connection
}

// Save save repository model in record
func (r *Repository) Save() error {
	result, err := r.db.Exec("insert into repositories (repository_id, owner, name, webhook_key, created_at) values (?, ?, ?, ?, now());", r.RepositoryID, r.Owner, r.Name, r.WebhookKey)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	r.ID, _ = result.LastInsertId()
	return nil
}

// FindByProjectID returns a repository related a project.
func FindByProjectID(projectID int64) (*Repository, error) {
	db := database.SharedInstance().Connection
	rows, err := db.Query("select repositories.id, repositories.repository_id, repositories.owner, repositories.name, repositories.webhook_key from projects inner join repositories on repositories.id = projects.repository_id where projects.id = ?;", projectID)
	if err != nil {
		return nil, errors.Wrap(err, "find repository error")
	}

	var id, repositoryID int64
	var owner, name sql.NullString
	var webhookKey string
	for rows.Next() {
		err = rows.Scan(&id, &repositoryID, &owner, &name, &webhookKey)
		if err != nil {
			return nil, errors.Wrap(err, "find repository error")
		}
	}
	r := New(id, repositoryID, owner.String, name.String, webhookKey)
	return r, nil
}

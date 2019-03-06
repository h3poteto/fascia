package repo

import (
	"database/sql"

	"github.com/pkg/errors"
)

// Repository has repository record
type Repo struct {
	db *sql.DB
}

// New is build new Repository struct
func New(db *sql.DB) *Repo {
	return &Repo{
		db,
	}
}

// FindByGithubRepoID is return a Repository struct from repository_id
func (r *Repo) FindByGithubRepoID(githubRepoID int64) (int64, int64, sql.NullString, sql.NullString, string, error) {
	var id, repoID int64
	var owner, name sql.NullString
	var webhookKey string
	err := r.db.QueryRow("select id, repository_id, owner, name, webhook_key from repositories where repository_id = ?;", githubRepoID).Scan(&id, &repoID, &owner, &name, &webhookKey)
	if err != nil {
		return 0, 0, sql.NullString{}, sql.NullString{}, "", errors.Wrap(err, "repo repository")
	}
	return id, repoID, owner, name, webhookKey, nil
}

// FindByProjectID returns a repository related a project.
func (r *Repo) FindByProjectID(projectID int64) (int64, int64, sql.NullString, sql.NullString, string, error) {
	var id, repoID int64
	var owner, name sql.NullString
	var webhookKey string
	err := r.db.QueryRow("select repositories.id, repositories.repository_id, repositories.owner, repositories.name, repositories.webhook_key from projects inner join repositories on repositories.id = projects.repository_id where projects.id = ?;", projectID).Scan(&id, &repoID, &owner, &name, &webhookKey)
	if err != nil {
		return 0, 0, sql.NullString{}, sql.NullString{}, "", errors.Wrap(err, "repo repository")
	}
	return id, repoID, owner, name, webhookKey, nil
}

// Create save repository model in record
func (r *Repo) Create(repositoryID int64, owner, name sql.NullString, webhookKey string) (int64, error) {
	result, err := r.db.Exec("insert into repositories (repository_id, owner, name, webhook_key, created_at) values (?, ?, ?, ?, now());", repositoryID, owner, name, webhookKey)
	if err != nil {
		return 0, errors.Wrap(err, "repo repository")
	}
	id, _ := result.LastInsertId()
	return id, nil
}

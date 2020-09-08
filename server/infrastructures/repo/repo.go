package repo

import (
	"database/sql"

	"github.com/h3poteto/fascia/server/domains/repo"
	"github.com/pkg/errors"
)

// Repo has db connection.
type Repo struct {
	db *sql.DB
}

// New is build new Repo struct
func New(db *sql.DB) *Repo {
	return &Repo{
		db,
	}
}

// FindByGithubRepoID is return a Repository struct from repository_id
func (r *Repo) FindByGithubRepoID(githubRepoID int64) (*repo.Repo, error) {
	var id, repoID int64
	var owner, name sql.NullString
	var webhookKey string
	err := r.db.QueryRow("SELECT id, repository_id, owner, name, webhook_key FROM repositories WHERE repository_id = $1;", githubRepoID).Scan(&id, &repoID, &owner, &name, &webhookKey)
	if err != nil {
		return nil, errors.Wrap(err, "repo repository")
	}
	return &repo.Repo{
		ID:           id,
		RepositoryID: repoID,
		Owner:        owner,
		Name:         name,
		WebhookKey:   webhookKey,
	}, nil
}

// FindByProjectID returns a repository related a project.
func (r *Repo) FindByProjectID(projectID int64) (*repo.Repo, error) {
	var id, repoID int64
	var owner, name sql.NullString
	var webhookKey string
	rows, err := r.db.Query("SELECT repositories.id, repositories.repository_id, repositories.owner, repositories.name, repositories.webhook_key FROM projects INNER JOIN repositories ON repositories.id = projects.repository_id WHERE projects.id = $1;", projectID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(&id, &repoID, &owner, &name, &webhookKey)
		if err != nil {
			return nil, err
		}
	}
	if id == 0 || repoID == 0 {
		err := errors.New("Record not found")
		return nil, &repo.NotFoundError{Err: err}
	}
	return &repo.Repo{
		ID:           id,
		RepositoryID: repoID,
		Owner:        owner,
		Name:         name,
		WebhookKey:   webhookKey,
	}, nil
}

// Create save repository model in record
func (r *Repo) Create(repositoryID int64, owner, name sql.NullString, webhookKey string) (int64, error) {
	var id int64
	err := r.db.QueryRow("INSERT INTO repositories (repository_id, owner, name, webhook_key) VALUES ($1, $2, $3, $4) RETURNING id;", repositoryID, owner, name, webhookKey).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "repo repository")
	}
	return id, nil
}

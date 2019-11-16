package project

import (
	"database/sql"

	"github.com/h3poteto/fascia/server/domains/project"
	"github.com/pkg/errors"
)

// Project has project record
type Project struct {
	db *sql.DB
}

// New returns a new project object
func New(db *sql.DB) *Project {
	return &Project{
		db,
	}
}

// Find search a project according to id
func (p *Project) Find(projectID int64) (*project.Project, error) {
	var id, userID int64
	var repositoryID sql.NullInt64
	var title, description string
	var showIssues, showPullRequests bool
	err := p.db.QueryRow("SELECT id, user_id, repository_id, title, description, show_issues, show_pull_requests FROM projects WHERE id = $1;", projectID).Scan(&id, &userID, &repositoryID, &title, &description, &showIssues, &showPullRequests)
	if err != nil {
		return nil, err
	}
	return &project.Project{
		ID:               id,
		UserID:           userID,
		Title:            title,
		Description:      description,
		RepositoryID:     repositoryID,
		ShowIssues:       showIssues,
		ShowPullRequests: showPullRequests,
	}, nil
}

// FindByRepositoryID search projects according to repository id
func (p *Project) FindByRepositoryID(repoID int64) ([]*project.Project, error) {
	result := []*project.Project{}
	rows, err := p.db.Query("SELECT id, user_id, repository_id, title, description, show_issues, show_pull_requests FROM projects WHERE repository_id = $1;", repoID)
	if err != nil {
		return nil, errors.Wrap(err, "project repository")
	}
	for rows.Next() {
		var id, userID int64
		var repositoryID sql.NullInt64
		var title, description string
		var showIssues, showPullRequests bool
		err = rows.Scan(&id, &userID, &repositoryID, &title, &description, &showIssues, &showPullRequests)
		if err != nil {
			return nil, errors.Wrap(err, "project repository")
		}
		p := &project.Project{
			ID:               id,
			UserID:           userID,
			Title:            title,
			Description:      description,
			RepositoryID:     repositoryID,
			ShowIssues:       showIssues,
			ShowPullRequests: showPullRequests,
		}
		result = append(result, p)
	}
	return result, nil
}

// Create save project model in record
func (p *Project) Create(userID int64, title string, description string, repositoryID sql.NullInt64, showIssues bool, showPullRequests bool, tx *sql.Tx) (int64, error) {
	var err error
	var id int64
	if tx != nil {
		err = tx.QueryRow("INSERT INTO projects (user_id, repository_id, title, description, show_issues, show_pull_requests) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;", userID, repositoryID, title, description, showIssues, showPullRequests).Scan(&id)
	} else {
		err = p.db.QueryRow("INSERT INTO projects (user_id, repository_id, title, description, show_issues, show_pull_requests) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;", userID, repositoryID, title, description, showIssues, showPullRequests).Scan(&id)
	}
	if err != nil {
		return 0, errors.Wrap(err, "project repository")
	}
	return id, nil
}

// Update update project model in record
func (p *Project) Update(id, userID int64, title string, description string, repositoryID sql.NullInt64, showIssues bool, showPullRequests bool) error {
	_, err := p.db.Exec("UPDATE projects SET user_id = $1, title = $2, description = $3, repository_id = $4, show_issues = $5, show_pull_requests = $6 WHERE id = $7;", userID, title, description, repositoryID, showIssues, showPullRequests, id)
	if err != nil {
		return errors.Wrap(err, "project repository")
	}
	return nil
}

// Delete delete a project model in record
func (p *Project) Delete(id int64) error {
	_, err := p.db.Exec("DELETE FROM projects WHERE id = $1;", id)
	if err != nil {
		return errors.Wrap(err, "project repository")
	}
	return nil
}

// Projects returns a project related a user.
func (p *Project) Projects(targetUserID int64) ([]*project.Project, error) {
	result := []*project.Project{}
	rows, err := p.db.Query("SELECT id, user_id, repository_id, title, description, show_issues, show_pull_requests FROM projects WHERE user_id = $1;", targetUserID)
	if err != nil {
		return nil, errors.Wrap(err, "project repository")
	}
	for rows.Next() {
		var id, userID int64
		var repositoryID sql.NullInt64
		var title, description string
		var showIssues, showPullRequests bool
		err := rows.Scan(&id, &userID, &repositoryID, &title, &description, &showIssues, &showPullRequests)
		if err != nil {
			return nil, errors.Wrap(err, "project repository")
		}
		p := &project.Project{
			ID:               id,
			UserID:           userID,
			Title:            title,
			Description:      description,
			RepositoryID:     repositoryID,
			ShowIssues:       showIssues,
			ShowPullRequests: showPullRequests,
		}
		result = append(result, p)
	}
	return result, nil
}

// OauthToken get oauth token related this project
func (p *Project) OauthToken(id int64) (string, error) {
	var oauthToken sql.NullString
	err := p.db.QueryRow("SELECT users.oauth_token FROM projects LEFT JOIN users ON users.id = projects.user_id WHERE projects.id = $1;", id).Scan(&oauthToken)
	if err != nil {
		return "", errors.Wrap(err, "project repository")
	}
	if !oauthToken.Valid {
		return "", errors.New("oauth token isn't exist")
	}

	return oauthToken.String, nil
}

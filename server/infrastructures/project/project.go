package project

import (
	"database/sql"

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
func (p *Project) Find(projectID int64) (int64, int64, string, string, sql.NullInt64, bool, bool, error) {
	var id, userID int64
	var repositoryID sql.NullInt64
	var title, description string
	var showIssues, showPullRequests bool
	err := p.db.QueryRow("select id, user_id, repository_id, title, description, show_issues, show_pull_requests from projects where id = ?;", projectID).Scan(&id, &userID, &repositoryID, &title, &description, &showIssues, &showPullRequests)
	return id, userID, title, description, repositoryID, showIssues, showPullRequests, err
}

// FindByRepositoryID search projects according to repository id
func (p *Project) FindByRepositoryID(repoID int64) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	rows, err := p.db.Query("select id, user_id, repository_id, title, description, show_issues, show_pull_requests from projects where repository_id = ?;", repoID)
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
		p := map[string]interface{}{
			"id":               id,
			"userID":           userID,
			"title":            title,
			"description":      description,
			"repositoryID":     repositoryID,
			"showIssues":       showIssues,
			"showPullRequests": showPullRequests,
		}
		result = append(result, p)
	}
	return result, nil
}

// Create save project model in record
func (p *Project) Create(userID int64, title string, description string, repositoryID sql.NullInt64, showIssues bool, showPullRequests bool, tx *sql.Tx) (int64, error) {
	var err error
	var result sql.Result
	if tx != nil {
		result, err = tx.Exec("insert into projects (user_id, repository_id, title, description, show_issues, show_pull_requests, created_at) values (?, ?, ?, ?, ?, ?, now());", userID, repositoryID, title, description, showIssues, showPullRequests)
	} else {
		result, err = p.db.Exec("insert into projects (user_id, repository_id, title, description, show_issues, show_pull_requests, created_at) values (?, ?, ?, ?, ?, ?, now());", userID, repositoryID, title, description, showIssues, showPullRequests)
	}
	if err != nil {
		return 0, errors.Wrap(err, "project repository")
	}
	id, _ := result.LastInsertId()
	return id, nil
}

// Update update project model in record
func (p *Project) Update(id, userID int64, title string, description string, repositoryID sql.NullInt64, showIssues bool, showPullRequests bool) error {
	_, err := p.db.Exec("update projects set user_id = ?, title = ?, description = ?, repository_id = ?, show_issues = ?, show_pull_requests = ? where id = ?;", userID, title, description, repositoryID, showIssues, showPullRequests, id)
	if err != nil {
		return errors.Wrap(err, "project repository")
	}
	return nil
}

// Delete delete a project model in record
func (p *Project) Delete(id int64) error {
	_, err := p.db.Exec("DELETE FROM projects WHERE id = ?;", id)
	if err != nil {
		return errors.Wrap(err, "project repository")
	}
	return nil
}

// Projects returns a project related a user.
func (p *Project) Projects(targetUserID int64) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	rows, err := p.db.Query("select id, user_id, repository_id, title, description, show_issues, show_pull_requests from projects where user_id = ?;", targetUserID)
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
		p := map[string]interface{}{
			"id":               id,
			"userID":           userID,
			"title":            title,
			"description":      description,
			"repositoryID":     repositoryID,
			"showIssues":       showIssues,
			"showPullRequests": showPullRequests,
		}
		result = append(result, p)
	}
	return result, nil
}

// OauthToken get oauth token related this project
func (p *Project) OauthToken(id int64) (string, error) {
	var oauthToken sql.NullString
	err := p.db.QueryRow("select users.oauth_token from projects left join users on users.id = projects.user_id where projects.id = ?;", id).Scan(&oauthToken)
	if err != nil {
		return "", errors.Wrap(err, "project repository")
	}
	if !oauthToken.Valid {
		return "", errors.New("oauth token isn't exist")
	}

	return oauthToken.String, nil
}

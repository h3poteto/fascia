package project

import (
	"github.com/h3poteto/fascia/server/models/db"

	"database/sql"

	"github.com/pkg/errors"
)

// Project has project record
type Project struct {
	ID               int64
	UserID           int64
	Title            string
	Description      string
	RepositoryID     sql.NullInt64
	ShowIssues       bool
	ShowPullRequests bool
	database         *sql.DB
}

// New returns a new project object
func New(id int64, userID int64, title string, description string, repositoryID sql.NullInt64, showIssues bool, showPullRequests bool) *Project {
	if userID == 0 {
		return nil
	}

	project := &Project{ID: id, UserID: userID, Title: title, Description: description, RepositoryID: repositoryID, ShowIssues: showIssues, ShowPullRequests: showPullRequests}
	project.initialize()
	return project
}

// Find search a project according to id
func Find(projectID int64) (*Project, error) {
	database := db.SharedInstance().Connection

	var id, userID int64
	var repositoryID sql.NullInt64
	var title string
	var description string
	var showIssues, showPullRequests bool
	err := database.QueryRow("select id, user_id, repository_id, title, description, show_issues, show_pull_requests from projects where id = ?;", projectID).Scan(&id, &userID, &repositoryID, &title, &description, &showIssues, &showPullRequests)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	project := New(id, userID, title, description, repositoryID, showIssues, showPullRequests)
	return project, nil
}

// FindByRepositoryID search projects according to repository id
func FindByRepositoryID(repoID int64) ([]*Project, error) {
	database := db.SharedInstance().Connection

	var slice []*Project
	rows, err := database.Query("select id, user_id, repository_id, title, description, show_issues, show_pull_requests from projects where repository_id = ?;", repoID)
	if err != nil {
		return nil, errors.Wrap(err, "find project error")
	}
	for rows.Next() {
		var id, userID int64
		var repositoryID sql.NullInt64
		var title string
		var description string
		var showIssues, showPullRequests bool
		err = rows.Scan(&id, &userID, &repositoryID, &title, &description, &showIssues, &showPullRequests)
		if err != nil {
			return nil, errors.Wrap(err, "scan project error")
		}
		p := New(id, userID, title, description, repositoryID, showIssues, showPullRequests)
		slice = append(slice, p)
	}
	return slice, nil
}

func (p *Project) initialize() {
	p.database = db.SharedInstance().Connection
}

// Save save project model in record
func (p *Project) Save(tx *sql.Tx) error {
	var err error
	var result sql.Result
	if tx != nil {
		result, err = tx.Exec("insert into projects (user_id, repository_id, title, description, show_issues, show_pull_requests, created_at) values (?, ?, ?, ?, ?, ?, now());", p.UserID, p.RepositoryID, p.Title, p.Description, p.ShowIssues, p.ShowPullRequests)
	} else {
		result, err = p.database.Exec("insert into projects (user_id, repository_id, title, description, show_issues, show_pull_requests, created_at) values (?, ?, ?, ?, ?, ?, now());", p.UserID, p.RepositoryID, p.Title, p.Description, p.ShowIssues, p.ShowPullRequests)
	}
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	p.ID, _ = result.LastInsertId()
	return nil
}

// Update update project model in record
func (p *Project) Update(title string, description string, showIssues bool, showPullRequests bool) error {
	_, err := p.database.Exec("update projects set title = ?, description = ?, show_issues = ?, show_pull_requests = ? where id = ?;", title, description, showIssues, showPullRequests, p.ID)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	p.Title = title
	p.Description = description
	p.ShowIssues = showIssues
	p.ShowPullRequests = showPullRequests

	return nil
}

// Delete delete a project model in record
func (p *Project) Delete() error {
	_, err := p.database.Exec("DELETE FROM projects WHERE id = ?;", p.ID)
	if err != nil {
		return err
	}
	return nil
}

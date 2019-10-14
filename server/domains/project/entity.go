package project

import (
	"database/sql"
)

// Project has a project model object
type Project struct {
	ID               int64
	UserID           int64
	Title            string
	Description      string
	RepositoryID     sql.NullInt64
	ShowIssues       bool
	ShowPullRequests bool
}

// New returns a project entity
func New(id int64, userID int64, title string, description string, repositoryID sql.NullInt64, showIssues bool, showPullRequests bool) *Project {
	return &Project{
		id,
		userID,
		title,
		description,
		repositoryID,
		showIssues,
		showPullRequests,
	}
}

// Update call project model update
func (p *Project) Update(title string, description string, showIssues bool, showPullRequests bool) error {
	p.Title = title
	p.Description = description
	p.ShowIssues = showIssues
	p.ShowPullRequests = showPullRequests
	return nil
}

// CheckOwner returns either owner of the project.
func (p *Project) CheckOwner(userID int64) bool {
	return p.UserID == userID
}

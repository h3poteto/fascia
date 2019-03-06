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
	infrastructure   Repository
}

type Repository interface {
	Find(int64) (int64, int64, string, string, sql.NullInt64, bool, bool, error)
	FindByRepositoryID(int64) ([]map[string]interface{}, error)
	Create(int64, string, string, sql.NullInt64, bool, bool, *sql.Tx) (int64, error)
	Update(int64, int64, string, string, sql.NullInt64, bool, bool) error
	Delete(int64) error
	Projects(int64) ([]map[string]interface{}, error)
	OauthToken(int64) (string, error)
}

// New returns a project entity
func New(id int64, userID int64, title string, description string, repositoryID sql.NullInt64, showIssues bool, showPullRequests bool, infrastructure Repository) *Project {
	return &Project{
		id,
		userID,
		title,
		description,
		repositoryID,
		showIssues,
		showPullRequests,
		infrastructure,
	}
}

// Create call project model save
func (p *Project) Create(tx *sql.Tx) error {
	id, err := p.infrastructure.Create(p.UserID, p.Title, p.Description, p.RepositoryID, p.ShowIssues, p.ShowPullRequests, tx)
	if err != nil {
		return err
	}
	p.ID = id
	return nil
}

// Update call project model update
func (p *Project) Update(title string, description string, showIssues bool, showPullRequests bool) error {
	if err := p.infrastructure.Update(p.ID, p.UserID, title, description, p.RepositoryID, showIssues, showPullRequests); err != nil {
		return err
	}
	p.Title = title
	p.Description = description
	p.ShowIssues = showIssues
	p.ShowPullRequests = showPullRequests
	return nil
}

// Delete delete a project model
func (p *Project) Delete() error {
	return p.infrastructure.Delete(p.ID)
}

// OauthToken get oauth token related this project
func (p *Project) OauthToken() (string, error) {
	return p.infrastructure.OauthToken(p.ID)
}

func (p *Project) CheckOwner(userID int64) bool {
	return p.UserID == userID
}

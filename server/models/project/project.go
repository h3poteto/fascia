package project

import (
	"github.com/h3poteto/fascia/server/models/db"

	"database/sql"

	"github.com/pkg/errors"
)

type ProjectStruct struct {
	ID               int64
	UserID           int64
	Title            string
	Description      string
	RepositoryID     sql.NullInt64
	ShowIssues       bool
	ShowPullRequests bool
	database         *sql.DB
}

func New(id int64, userID int64, title string, description string, repositoryID sql.NullInt64, showIssues bool, showPullRequests bool) *ProjectStruct {
	if userID == 0 {
		return nil
	}

	project := &ProjectStruct{ID: id, UserID: userID, Title: title, Description: description, RepositoryID: repositoryID, ShowIssues: showIssues, ShowPullRequests: showPullRequests}
	project.initialize()
	return project
}

func Find(projectID int64) (*ProjectStruct, error) {
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

func FindByRepositoryID(repoID int64) (*ProjectStruct, error) {
	database := db.SharedInstance().Connection

	var id, userID int64
	var repositoryID sql.NullInt64
	var title string
	var description string
	var showIssues, showPullRequests bool
	err := database.QueryRow("select id, user_id, repository_id, title, description, show_issues, show_pull_requests from projects where repository_id = ?;", repositoryID).Scan(&id, &userID, &repositoryID, &title, &description, &showIssues, &showPullRequests)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	project := New(id, userID, title, description, repositoryID, showIssues, showPullRequests)
	return project, nil
}

func (u *ProjectStruct) initialize() {
	u.database = db.SharedInstance().Connection
}

func (u *ProjectStruct) Save(tx *sql.Tx) error {
	var err error
	var result sql.Result
	if tx != nil {
		result, err = tx.Exec("insert into projects (user_id, repository_id, title, description, show_issues, show_pull_requests, created_at) values (?, ?, ?, ?, ?, ?, now());", u.UserID, u.RepositoryID, u.Title, u.Description, u.ShowIssues, u.ShowPullRequests)
	} else {
		result, err = u.database.Exec("insert into projects (user_id, repository_id, title, description, show_issues, show_pull_requests, created_at) values (?, ?, ?, ?, ?, ?, now());", u.UserID, u.RepositoryID, u.Title, u.Description, u.ShowIssues, u.ShowPullRequests)
	}
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	u.ID, _ = result.LastInsertId()
	return nil
}

func (u *ProjectStruct) Update(title string, description string, showIssues bool, showPullRequests bool) error {
	_, err := u.database.Exec("update projects set title = ?, description = ?, show_issues = ?, show_pull_requests = ? where id = ?;", title, description, showIssues, showPullRequests, u.ID)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	u.Title = title
	u.Description = description
	u.ShowIssues = showIssues
	u.ShowPullRequests = showPullRequests

	return nil
}

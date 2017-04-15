package views

import (
	"github.com/h3poteto/fascia/server/entities/project"
)

// Project provides a response structure for project
type Project struct {
	ID               int64  `json:ID`
	UserID           int64  `json:UserID`
	Title            string `json:Title`
	Description      string `json:Description`
	ShowIssues       bool   `json:ShowIssues`
	ShowPullRequests bool   `json:ShowPullRequests`
	RepositoryID     int64  `json:RepositoryID`
}

// ParseProjectJSON returns a project struct for response
func ParseProjectJSON(project *project.Project) (*Project, error) {
	var repositoryID int64
	repo, find, err := project.Repository()
	if err != nil {
		return nil, err
	}
	if find {
		repositoryID = repo.RepositoryModel.ID
	}

	return &Project{
		ID:               project.ProjectModel.ID,
		UserID:           project.ProjectModel.UserID,
		Title:            project.ProjectModel.Title,
		Description:      project.ProjectModel.Description,
		ShowIssues:       project.ProjectModel.ShowIssues,
		ShowPullRequests: project.ProjectModel.ShowPullRequests,
		RepositoryID:     repositoryID,
	}, nil
}

// ParseProjectsJSON returns some projects structs for response
func ParseProjectsJSON(projects []*project.Project) ([]*Project, error) {
	var results []*Project
	for _, p := range projects {
		parse, err := ParseProjectJSON(p)
		if err != nil {
			return nil, err
		}
		results = append(results, parse)
	}
	return results, nil
}

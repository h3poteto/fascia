package views

import (
	"github.com/h3poteto/fascia/server/domains/project"
	"github.com/h3poteto/fascia/server/usecases/board"
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
	repo, err := board.ProjectRepository(project)
	if err == nil {
		repositoryID = repo.ID
	}

	return &Project{
		ID:               project.ID,
		UserID:           project.UserID,
		Title:            project.Title,
		Description:      project.Description,
		ShowIssues:       project.ShowIssues,
		ShowPullRequests: project.ShowPullRequests,
		RepositoryID:     repositoryID,
	}, nil
}

// ParseProjectsJSON returns some projects structs for response
func ParseProjectsJSON(projects []*project.Project) ([]*Project, error) {
	results := []*Project{}
	for _, p := range projects {
		parse, err := ParseProjectJSON(p)
		if err != nil {
			return nil, err
		}
		results = append(results, parse)
	}
	return results, nil
}

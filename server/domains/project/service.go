package project

import (
	"database/sql"

	"github.com/h3poteto/fascia/server/domains/repo"
	"github.com/pkg/errors"
)

// Find returns a project entity
func Find(targetID int64, infrastructure Repository) (*Project, error) {
	id, userID, title, description, repositoryID, showIssues, showPullRequests, err := infrastructure.Find(targetID)
	if err != nil {
		return nil, err
	}
	return New(id, userID, title, description, repositoryID, showIssues, showPullRequests, infrastructure), nil
}

// FindByRepoID returns project entities.
func FindByRepoID(targetRepoID int64, infrastructure Repository) ([]*Project, error) {
	projects, err := infrastructure.FindByRepositoryID(targetRepoID)
	if err != nil {
		return nil, err
	}
	var result []*Project
	for _, project := range projects {
		id, ok := project["id"].(int64)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		userID, ok := project["userID"].(int64)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		title, ok := project["title"].(string)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		description, ok := project["description"].(string)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		repositoryID, ok := project["repositoryID"].(sql.NullInt64)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		showIssues, ok := project["showIssues"].(bool)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		showPullRequests, ok := project["showPullRequests"].(bool)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		p := New(id, userID, title, description, repositoryID, showIssues, showPullRequests, infrastructure)
		result = append(result, p)

	}
	return result, nil
}

// Projects returns all projects related a user.
func Projects(targetUserID int64, infrastructure Repository) ([]*Project, error) {
	var result []*Project

	projects, err := infrastructure.Projects(targetUserID)
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		id, ok := project["id"].(int64)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		userID, ok := project["userID"].(int64)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		title, ok := project["title"].(string)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		description, ok := project["description"].(string)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		repositoryID, ok := project["repositoryID"].(sql.NullInt64)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		showIssues, ok := project["showIssues"].(bool)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		showPullRequests, ok := project["showPullRequests"].(bool)
		if !ok {
			return nil, errors.New("Can not convert interface")
		}
		p := New(id, userID, title, description, repositoryID, showIssues, showPullRequests, infrastructure)
		result = append(result, p)

	}
	return result, nil
}

// Repository returns a repository entity related this project
func (p *Project) Repository(infrastructure repo.Repository) (*repo.Repo, error) {
	return repo.FindByProjectID(p.ID, infrastructure)
}

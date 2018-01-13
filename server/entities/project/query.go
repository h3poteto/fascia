package project

import (
	"github.com/h3poteto/fascia/server/entities/list"
	"github.com/h3poteto/fascia/server/entities/repository"
	"github.com/h3poteto/fascia/server/infrastructures/project"
)

// Find returns a project entity
func Find(id int64) (*Project, error) {
	p := &Project{
		ID: id,
	}
	if err := p.reload(); err != nil {
		return nil, err
	}
	return p, nil
}

// FindByRepositoryID returns project entities
func FindByRepositoryID(repositoryID int64) ([]*Project, error) {
	projects, err := project.FindByRepositoryID(repositoryID)
	if err != nil {
		return nil, err
	}
	var slice []*Project
	for _, p := range projects {
		s := &Project{
			infrastructure: p,
		}
		if err := s.reload(); err != nil {
			return nil, err
		}
		slice = append(slice, s)
	}
	return slice, nil
}

// Projects returns all projects related a user.
func Projects(userID int64) ([]*Project, error) {
	var slice []*Project

	projects, err := project.Projects(userID)
	if err != nil {
		return nil, err
	}

	for _, p := range projects {
		s := &Project{
			infrastructure: p,
		}
		if err := s.reload(); err != nil {
			return nil, err
		}
		slice = append(slice, s)
	}
	return slice, nil
}

// OauthToken get oauth token related this project
func (p *Project) OauthToken() (string, error) {
	return p.infrastructure.OauthToken()
}

// Lists list up lists related this project
func (p *Project) Lists() ([]*list.List, error) {
	return list.Lists(p.ID)
}

// NoneList returns a none list related this project
func (p *Project) NoneList() (*list.List, error) {
	// noneが存在しないということはProjectsController#Createがうまく行ってないので，そっちでエラーハンドリングしてほしい
	return list.NoneList(p.ID)
}

// Repository returns a repository entity related this project
// If repository does not exist, return false
func (p *Project) Repository() (*repository.Repository, bool, error) {
	return repository.FindByProjectID(p.ID)
}

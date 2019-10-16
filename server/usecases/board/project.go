package board

import (
	"database/sql"

	"github.com/h3poteto/fascia/server/domains/list"
	domain "github.com/h3poteto/fascia/server/domains/project"
	"github.com/h3poteto/fascia/server/domains/repo"
	"github.com/h3poteto/fascia/server/domains/services"
	repository "github.com/h3poteto/fascia/server/infrastructures/project"
)

// InjectProjectRepository returns a project Repository.
func InjectProjectRepository() domain.Repository {
	return repository.New(InjectDB())
}

// FindProject finds a project.
func FindProject(id int64) (*domain.Project, error) {
	infra := InjectProjectRepository()
	return infra.Find(id)
}

// CreateProject create a project and sync it to github.
func CreateProject(userID int64, title string, description string, repositoryID int64, oauthToken sql.NullString) (*domain.Project, error) {
	var repoID sql.NullInt64
	if repositoryID != 0 && oauthToken.Valid {
		infra := InjectRepoRepository()
		r, err := infra.FindByGithubRepoID(repositoryID)
		if err != nil {
			r, err = services.CreateRepo(repositoryID, oauthToken.String, InjectRepoRepository())
			if err != nil {
				return nil, err
			}
		}
		repoID = sql.NullInt64{Int64: r.ID, Valid: true}
	}

	tx, err := InjectDB().Begin()
	if err != nil {
		return nil, err
	}
	project := domain.New(0, userID, title, description, repoID, true, true)
	infra := InjectProjectRepository()
	id, err := infra.Create(project.UserID, project.Title, project.Description, project.RepositoryID, project.ShowIssues, project.ShowPullRequests, tx)
	if err != nil {
		return nil, err
	}
	project.ID = id
	err = services.CreateInitialLists(project, InjectListRepository(), tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	go services.AfterCreateProject(project, InjectProjectRepository(), InjectListRepository(), InjectTaskRepository(), InjectRepoRepository())

	return project, nil
}

// DeleteProject delete project and delete webhook
func DeleteProject(projectID int64) error {
	projectInfra := InjectProjectRepository()
	project, err := projectInfra.Find(projectID)
	if err != nil {
		return err
	}

	repoInfra := InjectRepoRepository()
	r, err := repoInfra.FindByProjectID(project.ID)
	if err == nil {
		token, _ := projectInfra.OauthToken(project.ID)
		r.DeleteWebhook(token)
	}
	err = services.DeleteLists(project, InjectListRepository())
	if err != nil {
		return err
	}
	return projectInfra.Delete(project.ID)
}

func ProjectRepository(p *domain.Project) (*repo.Repo, error) {
	repoInfra := InjectRepoRepository()
	return repoInfra.FindByProjectID(p.ID)
}

// ProjectLists returns all lists related the project.
func ProjectLists(p *domain.Project) ([]*list.List, error) {
	repo := InjectListRepository()
	return repo.Lists(p.ID)
}

// ProjectNoneList returns the none list related the project.
func ProjectNoneList(p *domain.Project) (*list.List, error) {
	repo := InjectListRepository()
	return repo.NoneList(p.ID)
}

// UpdateProject updates a project.
func UpdateProject(p *domain.Project, title, description string, showIssues, showPullRequests bool) error {
	p.Update(title, description, showIssues, showPullRequests)
	infra := InjectProjectRepository()
	return infra.Update(p.ID, p.UserID, p.Title, p.Description, p.RepositoryID, p.ShowIssues, p.ShowPullRequests)
}

// OauthTokenFromProject gets oauth token form specified project.
func OauthTokenFromProject(p *domain.Project) (string, error) {
	infra := InjectProjectRepository()
	return infra.OauthToken(p.ID)
}

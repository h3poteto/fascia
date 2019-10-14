package board

import (
	"database/sql"

	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/domains/list"
	domain "github.com/h3poteto/fascia/server/domains/project"
	"github.com/h3poteto/fascia/server/domains/repo"
	repository "github.com/h3poteto/fascia/server/infrastructures/project"
)

// InjectProjectRepository returns a project Repository.
func InjectProjectRepository() domain.Repository {
	return repository.New(InjectDB())
}

// FindProject finds a project.
func FindProject(id int64) (*domain.Project, error) {
	return domain.Find(id, InjectProjectRepository())
}

func findProjectByRepoID(repoID int64) ([]*domain.Project, error) {
	return domain.FindByRepoID(repoID, InjectProjectRepository())
}

// fetchCreatedInitialList fetch initial list to github
func fetchCreatedInitialList(p *domain.Project) error {
	repo, err := p.Repository(InjectRepoRepository())
	if err != nil {
		return err
	}

	oauthToken, err := p.OauthToken()
	if err != nil {
		return err
	}

	listRepo := InjectListRepository()
	lists, err := listRepo.Lists(p.ID)
	if err != nil {
		return err
	}
	for _, l := range lists {
		label, err := repo.CheckLabelPresent(oauthToken, l.Title.String)
		if err != nil {
			return err
		}
		if label != nil {
			continue
		}
		_, err = repo.CreateGithubLabel(oauthToken, l.Title.String, l.Color.String)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateProject create a project and sync it to github.
func CreateProject(userID int64, title string, description string, repositoryID int64, oauthToken sql.NullString) (*domain.Project, error) {
	var repoID sql.NullInt64
	if repositoryID != 0 && oauthToken.Valid {
		r, err := repo.FindByGithubRepoID(repositoryID, InjectRepoRepository())
		if err != nil {
			r, err = repo.CreateRepo(repositoryID, oauthToken.String, InjectRepoRepository())
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
	project := domain.New(0, userID, title, description, repoID, true, true, InjectProjectRepository())
	if err := project.Create(tx); err != nil {
		return nil, err
	}
	err = createInitialLists(project, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	go func(project *domain.Project) {
		// Create initial list before get issues from github
		err := fetchCreatedInitialList(project)
		if err != nil {
			return
		}
		// Sync issues from github
		_, err = FetchGithub(project)
		if err != nil {
			logging.SharedInstance().MethodInfoWithStacktrace("Project", "Create", err).Error(err)
			return
		}

		// Create Webhook in github
		r, err := project.Repository(InjectRepoRepository())
		if err != nil {
			logging.SharedInstance().MethodInfoWithStacktrace("Project", "Create", err).Error(err)
			return
		}
		token, err := project.OauthToken()
		if err != nil {
			logging.SharedInstance().MethodInfoWithStacktrace("Project", "Create", err).Error(err)
			return
		}
		err = r.CreateWebhook(token)
		if err != nil {
			logging.SharedInstance().MethodInfo("Project", "Create").Infof("failed to create webhook: %v", err)
			return
		}
		logging.SharedInstance().MethodInfo("Project", "Create").Info("success to create webhook")
	}(project)

	return project, nil
}

// createInitialLists create initial lists in self project
func createInitialLists(project *domain.Project, tx *sql.Tx) error {
	closeListOption, err := FindListOptionByAction("close")
	if err != nil {
		return err
	}
	todoName := config.Element("init_list").(map[interface{}]interface{})["todo"].(string)
	todo := list.New(
		0,
		project.ID,
		project.UserID,
		sql.NullString{String: todoName, Valid: true},
		sql.NullString{String: "f37b1d", Valid: true},
		false,
		nil,
	)
	inprogressName := config.Element("init_list").(map[interface{}]interface{})["inprogress"].(string)
	inprogress := list.New(
		0,
		project.ID,
		project.UserID,
		sql.NullString{String: inprogressName, Valid: true},
		sql.NullString{String: "5eb95e", Valid: true},
		false,
		nil,
	)
	doneName := config.Element("init_list").(map[interface{}]interface{})["done"].(string)
	done := list.New(
		0,
		project.ID,
		project.UserID,
		sql.NullString{String: doneName, Valid: true},
		sql.NullString{String: "333333", Valid: true},
		false,
		closeListOption,
	)
	noneName := config.Element("init_list").(map[interface{}]interface{})["none"].(string)
	none := list.New(
		0,
		project.ID,
		project.UserID,
		sql.NullString{String: noneName, Valid: true},
		sql.NullString{String: "ffffff", Valid: true},
		false,
		nil,
	)

	repo := InjectListRepository()
	// This method only save lists.
	// Use another methods to sync github.
	if _, err := repo.Create(none.ProjectID, none.UserID, none.Title, none.Color, sql.NullInt64{}, none.IsHidden, tx); err != nil {
		return err
	}

	if _, err := repo.Create(todo.ProjectID, todo.UserID, todo.Title, todo.Color, sql.NullInt64{}, todo.IsHidden, tx); err != nil {
		return err
	}
	if _, err := repo.Create(inprogress.ProjectID, inprogress.UserID, inprogress.Title, inprogress.Color, sql.NullInt64{}, inprogress.IsHidden, tx); err != nil {
		return err
	}
	if _, err := repo.Create(done.ProjectID, done.UserID, done.Title, done.Color, sql.NullInt64{Int64: done.Option.ID, Valid: true}, done.IsHidden, tx); err != nil {
		return err
	}
	return nil
}

// DeleteProject delete project and delete webhook
func DeleteProject(projectID int64) error {
	project, err := FindProject(projectID)
	if err != nil {
		return err
	}

	r, err := project.Repository(InjectRepoRepository())
	if err == nil {
		token, _ := project.OauthToken()
		r.DeleteWebhook(token)
	}
	err = deleteLists(project)
	if err != nil {
		return err
	}
	return project.Delete()
}

// deleteLists delete all lists related a project
func deleteLists(p *domain.Project) error {
	repoList := InjectListRepository()
	lists, err := repoList.Lists(p.ID)
	if err != nil {
		return err
	}
	for _, l := range lists {
		err := repoList.DeleteTasks(l.ID)
		if err != nil {
			return err
		}
		err = repoList.Delete(l.ID)
		if err != nil {
			return err
		}
	}
	noneList, err := repoList.NoneList(p.ID)
	err = repoList.DeleteTasks(noneList.ID)
	if err != nil {
		return err
	}
	return repoList.Delete(noneList.ID)
}

// ProjectRepository returns a repo related the project.
func ProjectRepository(p *domain.Project) (*repo.Repo, error) {
	return p.Repository(InjectRepoRepository())
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

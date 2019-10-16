package services

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/hub"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/domains/list"
	"github.com/h3poteto/fascia/server/domains/project"
	"github.com/h3poteto/fascia/server/domains/repo"
	"github.com/h3poteto/fascia/server/domains/task"
)

// AfterCreateProject sync issues from github.
func AfterCreateProject(p *project.Project, projectInfra project.Repository, listInfra list.Repository, taskInfra task.Repository, repoInfra repo.Repository) {
	// Create initial list before get issues from github
	err := fetchCreatedInitialList(p, projectInfra, listInfra, repoInfra)
	if err != nil {
		return
	}
	// Sync issues from github
	_, err = FetchGithub(p, projectInfra, listInfra, taskInfra, repoInfra)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("Project", "Create", err).Error(err)
		return
	}

	// Create Webhook in github
	r, err := repoInfra.FindByProjectID(p.ID)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("Project", "Create", err).Error(err)
		return
	}
	token, err := projectInfra.OauthToken(p.ID)
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
	return
}

// fetchCreatedInitialList fetch initial list to github
func fetchCreatedInitialList(p *project.Project, projectInfra project.Repository, listInfra list.Repository, repoInfra repo.Repository) error {
	repo, err := repoInfra.FindByProjectID(p.ID)
	if err != nil {
		return err
	}

	oauthToken, err := projectInfra.OauthToken(p.ID)
	if err != nil {
		return err
	}

	lists, err := listInfra.Lists(p.ID)
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

// CreateInitialLists create initial lists in self project
func CreateInitialLists(project *project.Project, listInfra list.Repository, tx *sql.Tx) error {
	closeListOption, err := listInfra.FindOptionByAction("close")
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

	// This method only save lists.
	// Use another methods to sync github.
	if _, err := listInfra.Create(none.ProjectID, none.UserID, none.Title, none.Color, sql.NullInt64{}, none.IsHidden, tx); err != nil {
		return err
	}

	if _, err := listInfra.Create(todo.ProjectID, todo.UserID, todo.Title, todo.Color, sql.NullInt64{}, todo.IsHidden, tx); err != nil {
		return err
	}
	if _, err := listInfra.Create(inprogress.ProjectID, inprogress.UserID, inprogress.Title, inprogress.Color, sql.NullInt64{}, inprogress.IsHidden, tx); err != nil {
		return err
	}
	if _, err := listInfra.Create(done.ProjectID, done.UserID, done.Title, done.Color, sql.NullInt64{Int64: done.Option.ID, Valid: true}, done.IsHidden, tx); err != nil {
		return err
	}
	return nil
}

// CreateRepo create repository record based on github repository
func CreateRepo(targetRepositoryID int64, oauthToken string, repoInfra repo.Repository) (*repo.Repo, error) {
	// confirm github
	h := hub.New(oauthToken)
	githubRepo, err := h.GetRepository(int(targetRepositoryID))
	if err != nil {
		return nil, err
	}
	// generate webhook key
	key := generateWebhookKey(*githubRepo.Name)
	owner := sql.NullString{String: *githubRepo.Owner.Login, Valid: true}
	name := sql.NullString{String: *githubRepo.Name, Valid: true}
	_, err = repoInfra.Create(int64(*githubRepo.ID), owner, name, key)
	if err != nil {
		return nil, err
	}
	repo, err := repoInfra.FindByGithubRepoID(int64(*githubRepo.ID))
	return repo, nil
}

// generateWebhookKey create new md5 hash
func generateWebhookKey(seed string) string {
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	io.WriteString(h, seed)
	token := fmt.Sprintf("%x", h.Sum(nil))

	return token
}

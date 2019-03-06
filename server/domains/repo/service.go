package repo

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/hub"
)

// FindByGithubRepoID find repository entity according to repository id in github
func FindByGithubRepoID(targetRepositoryID int64, infrastructure Repository) (*Repo, error) {
	id, repositoryID, owner, name, webhookKey, err := infrastructure.FindByGithubRepoID(targetRepositoryID)
	if err != nil {
		return nil, err
	}
	return New(id, repositoryID, owner, name, webhookKey, infrastructure), nil
}

// FindByProjectID returns a repository related a project.
func FindByProjectID(targetProjectID int64, infrastructure Repository) (*Repo, error) {
	id, repositoryID, owner, name, webhookKey, err := infrastructure.FindByProjectID(targetProjectID)
	if err != nil {
		return nil, err
	}
	return New(id, repositoryID, owner, name, webhookKey, infrastructure), nil
}

// CreateRepository create repository record based on github repository
func CreateRepo(targetRepositoryID int64, oauthToken string, infrastructure Repository) (*Repo, error) {
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
	repo := New(0, int64(*githubRepo.ID), owner, name, key, infrastructure)
	err = repo.Create()
	if err != nil {
		return nil, err
	}
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

func (r *Repo) CreateWebhook(oauthToken string) error {
	url := fmt.Sprintf("%s://%s/repositories/hooks/github", config.Element("protocol").(string), config.Element("fqdn"))
	// If the webhook already exist, we update the webhook.
	hook, err := r.SearchWebhookInGithub(oauthToken, url)
	if err != nil {
		return err
	}
	if hook != nil {
		return r.UpdateWebhookInGithub(oauthToken, url, hook)
	}
	return r.CreateWebhookInGithub(oauthToken, url)
}

// DeleteWebhook call DeleteWebhook if project has repository
func (r *Repo) DeleteWebhook(oauthToken string) error {
	url := fmt.Sprintf("%s://%s/repositories/hooks/github", config.Element("protocol").(string), config.Element("fqdn"))
	hook, err := r.SearchWebhookInGithub(oauthToken, url)
	if err != nil {
		return err
	}
	if hook != nil {
		return r.DeleteWebhookInGithub(oauthToken, hook)
	}
	return nil
}

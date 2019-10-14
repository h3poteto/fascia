package repo

import (
	"crypto/hmac"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"

	"github.com/h3poteto/fascia/config"
	"github.com/pkg/errors"
)

// Repo has repository model object
type Repo struct {
	ID           int64
	RepositoryID int64
	Owner        sql.NullString
	Name         sql.NullString
	WebhookKey   string
}

// New returns a repository entity
func New(id int64, repositoryID int64, owner sql.NullString, name sql.NullString, webhookKey string) *Repo {
	return &Repo{
		id,
		repositoryID,
		owner,
		name,
		webhookKey,
	}
}

// Authenticate is check token and webhookKey with response
func (r *Repo) Authenticate(token string, response []byte) error {
	mac := hmac.New(sha1.New, []byte(r.WebhookKey))
	mac.Write(response)
	hashedToken := hex.EncodeToString(mac.Sum(nil))
	if token != ("sha1=" + hashedToken) {
		return errors.New("token is not equal webhookKey")
	}
	return nil
}

// CreateWebhook creates or updates the webhook.
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

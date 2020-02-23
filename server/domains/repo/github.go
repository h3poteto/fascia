package repo

import (
	"fmt"

	"github.com/google/go-github/github"
	"github.com/h3poteto/fascia/lib/modules/hub"
)

// CheckLabelPresent confirm existence label in github
func (r *Repo) CheckLabelPresent(token, title string) (*github.Label, error) {
	return hub.CheckLabelPresent(token, r.Owner.String, r.Name.String, title)
}

// CreateGithubLabel create new label in github
func (r *Repo) CreateGithubLabel(token, title, color string) (*github.Label, error) {
	return hub.CreateGithubLabel(token, r.Owner.String, r.Name.String, title, color)
}

// UpdateGithubLabel update exist label information in github
func (r *Repo) UpdateGithubLabel(token, originalTitle, title, color string) (*github.Label, error) {
	return hub.UpdateGithubLabel(token, r.Owner.String, r.Name.String, originalTitle, title, color)
}

// CreateGithubIssue create new issue in github
func (r *Repo) CreateGithubIssue(token, title, description string, labels []string) (*github.Issue, error) {
	return hub.CreateGithubIssue(token, r.Owner.String, r.Name.String, title, description, labels)
}

// EditGithubIssue update exist issue in github
func (r *Repo) EditGithubIssue(token, title, description, state string, issueNumber int, labels []string) (bool, error) {
	return hub.EditGithubIssue(token, r.Owner.String, r.Name.String, title, description, state, issueNumber, labels)
}

// GetGithubIssue return a issue in github
func (r *Repo) GetGithubIssue(token string, number int) (*github.Issue, error) {
	return hub.GetGithubIssue(token, r.Owner.String, r.Name.String, number)
}

// GetGithubIssues return few issues in github
func (r *Repo) GetGithubIssues(token string) ([]*github.Issue, []*github.Issue, error) {
	return hub.GetGithubIssues(token, r.Owner.String, r.Name.String)
}

// ListLabels returns all labels in github
func (r *Repo) ListLabels(token string) ([]*github.Label, error) {
	return hub.ListLabels(token, r.Owner.String, r.Name.String)
}

// CreateWebhookInGithub create a new webhook in github
func (r *Repo) CreateWebhookInGithub(token, url string) error {
	return hub.CreateWebhook(token, r.Owner.String, r.Name.String, r.WebhookKey, url)
}

// UpdateWebhookInGithub update a exist webhook in github
func (r *Repo) UpdateWebhookInGithub(token, url string, hook *github.Hook) error {
	return hub.EditWebhook(token, r.Owner.String, r.Name.String, r.WebhookKey, url, hook)
}

func (r *Repo) listWebhooksInGithub(token string) ([]*github.Hook, error) {
	return hub.ListWebhooks(token, r.Owner.String, r.Name.String)
}

// SearchWebhookInGithub search a webhook according to configured url
func (r *Repo) SearchWebhookInGithub(token, url string) (*github.Hook, error) {
	hooks, err := r.listWebhooksInGithub(token)
	if err != nil {
		return nil, err
	}
	for _, h := range hooks {
		config := h.Config
		fmt.Printf("debug: %+v", config)
		// Sometimes url does not exist in webhook config.
		if val, ok := config["url"]; ok {
			if val.(string) == url {
				return h, nil
			}
		}
	}
	return nil, nil
}

// DeleteWebhookInGithub delete a webhook
func (r *Repo) DeleteWebhookInGithub(token string, hook *github.Hook) error {
	return hub.DeleteWebhook(token, r.Owner.String, r.Name.String, hook)
}

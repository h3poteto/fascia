package repository

import (
	"github.com/google/go-github/github"
	"github.com/h3poteto/fascia/lib/modules/hub"
)

// CheckLabelPresent confirm existance label in github
func (r *Repository) CheckLabelPresent(token, title string) (*github.Label, error) {
	return hub.CheckLabelPresent(token, r.RepositoryModel.Owner.String, r.RepositoryModel.Name.String, title)
}

// CreateGithubLabel create new label in github
func (r *Repository) CreateGithubLabel(token, title, color string) (*github.Label, error) {
	return hub.CreateGithubLabel(token, r.RepositoryModel.Owner.String, r.RepositoryModel.Name.String, title, color)
}

// UpdateGithubLabel update exist label information in github
func (r *Repository) UpdateGithubLabel(token, originalTitle, title, color string) (*github.Label, error) {
	return hub.UpdateGithubLabel(token, r.RepositoryModel.Owner.String, r.RepositoryModel.Name.String, originalTitle, title, color)
}

// CreateGithubIssue create new issue in github
func (r *Repository) CreateGithubIssue(token, title, description string, labels []string) (*github.Issue, error) {
	return hub.CreateGithubIssue(token, r.RepositoryModel.Owner.String, r.RepositoryModel.Name.String, title, description, labels)
}

// EditGithubIssue update exist issue in github
func (r *Repository) EditGithubIssue(token, title, description, state string, issueNumber int, labels []string) (bool, error) {
	return hub.EditGithubIssue(token, r.RepositoryModel.Owner.String, r.RepositoryModel.Name.String, title, description, state, issueNumber, labels)
}

// GetGithubIssue return a issue in github
func (r *Repository) GetGithubIssue(token string, number int) (*github.Issue, error) {
	return hub.GetGithubIssue(token, r.RepositoryModel.Owner.String, r.RepositoryModel.Name.String, number)
}

// GetGithubIssues return few issues in github
func (r *Repository) GetGithubIssues(token string) ([]*github.Issue, []*github.Issue, error) {
	return hub.GetGithubIssues(token, r.RepositoryModel.Owner.String, r.RepositoryModel.Name.String)
}

// ListLabels returns all labels in github
func (r *Repository) ListLabels(token string) ([]*github.Label, error) {
	return hub.ListLabels(token, r.RepositoryModel.Owner.String, r.RepositoryModel.Name.String)
}

// CreateWebhook create a new webhook in github
func (r *Repository) CreateWebhook(token, url string) error {
	return hub.CreateWebhook(token, r.RepositoryModel.Owner.String, r.RepositoryModel.Name.String, r.RepositoryModel.WebhookKey, url)
}

// UpdateWebhook update a exist webhook in github
func (r *Repository) UpdateWebhook(token, url string, hook *github.Hook) error {
	return hub.EditWebhook(token, r.RepositoryModel.Owner.String, r.RepositoryModel.Name.String, r.RepositoryModel.WebhookKey, url, hook)
}

func (r *Repository) listWebhooks(token string) ([]*github.Hook, error) {
	return hub.ListWebhooks(token, r.RepositoryModel.Owner.String, r.RepositoryModel.Name.String)
}

// SearchWebhook search a webhook according to configured url
func (r *Repository) SearchWebhook(token, url string) (*github.Hook, error) {
	hooks, err := r.listWebhooks(token)
	if err != nil {
		return nil, err
	}
	for _, h := range hooks {
		config := h.Config
		if config["url"].(string) == url {
			return h, nil
		}
	}
	return nil, nil
}

// DeleteWebhook delete a webhook
func (r *Repository) DeleteWebhook(token string, hook *github.Hook) error {
	return hub.DeleteWebhook(token, r.RepositoryModel.Owner.String, r.RepositoryModel.Name.String, hook)
}

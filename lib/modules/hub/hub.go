package hub

import (
	"github.com/h3poteto/fascia/lib/modules/logging"

	"context"
	"time"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// Hub has github client struct
type Hub struct {
	client *github.Client
}

// New returns Hub struct
func New(token string) *Hub {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	return &Hub{client: client}
}

// AllRepositories returns all repositories in github account
func (h *Hub) AllRepositories() ([]*github.Repository, error) {
	nextPage := -1
	var repositories []*github.Repository
	for nextPage != 0 {
		if nextPage < 0 {
			nextPage = 0
		}
		repositoryOption := &github.RepositoryListOptions{
			Type:      "all",
			Sort:      "full_name",
			Direction: "asc",
			ListOptions: github.ListOptions{
				Page:    nextPage,
				PerPage: 50,
			},
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		repos, res, err := h.client.Repositories.List(ctx, "", repositoryOption)
		nextPage = res.NextPage
		if err != nil {
			err := errors.Wrap(err, "repository error")
			return nil, err
		}
		repositories = append(repositories, repos...)
	}
	return repositories, nil
}

// GetRepository returns a repository struct
func (h *Hub) GetRepository(ID int64) (*github.Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	repo, _, err := h.client.Repositories.GetByID(ctx, ID)
	if err != nil {
		return nil, errors.Wrap(err, "response error")
	}
	return repo, nil
}

// CheckLabelPresent confirm label existen in github
func CheckLabelPresent(token, owner, name, title string) (*github.Label, error) {
	client := prepareClient(token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	githubLabel, response, err := client.Issues.GetLabel(ctx, owner, name, title)
	logging.SharedInstance().MethodInfo("hub", "CheckLabelPresent").Debugf("respone of geting github label: %+v", response)
	if err != nil {
		logging.SharedInstance().MethodInfo("hub", "CheckLabelPresent").Debugf("cannot find github label: %v", name)
		// TODO: 本当はerrorで返したいが，error=nil, label=nilでラベルの存在判定をしている箇所があるので，それらを駆逐したい
		return nil, nil
	}
	logging.SharedInstance().MethodInfo("hub", "CheckLabelPresent").Debugf("github label is exist: %+v", githubLabel)
	return githubLabel, nil
}

// CreateGithubLabel create a new label in github
func CreateGithubLabel(token, owner, name, title, color string) (*github.Label, error) {
	client := prepareClient(token)

	label := &github.Label{
		Name:  &title,
		Color: &color,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	githubLabel, response, err := client.Issues.CreateLabel(ctx, owner, name, label)
	logging.SharedInstance().MethodInfo("hub", "CreateGithubLabel").Debugf("response of creating github label: %+v\n", response)
	if err != nil {
		return nil, errors.Wrap(err, "response is error")
	}
	logging.SharedInstance().MethodInfo("hub", "CreateGithubLabel").Debugf("github label is created: %+v", githubLabel)
	return githubLabel, nil
}

// UpdateGithubLabel update a exist label in github
func UpdateGithubLabel(token, owner, name, originalTitle, title, color string) (*github.Label, error) {
	client := prepareClient(token)

	label := &github.Label{
		Name:  &title,
		Color: &color,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	githubLabel, response, err := client.Issues.EditLabel(ctx, owner, name, originalTitle, label)
	logging.SharedInstance().MethodInfo("hub", "UpddateGithubLabel").Debugf("response of updating github label: %+v", response)
	if err != nil {
		return nil, errors.Wrap(err, "response is error")
	}
	logging.SharedInstance().MethodInfo("hub", "UpdateGithubLabel").Debugf("github label is updated: %+v", githubLabel)
	return githubLabel, nil
}

// CreateGithubIssue create a new issue in github
func CreateGithubIssue(token, owner, name, title, description string, labels []string) (*github.Issue, error) {
	client := prepareClient(token)

	issueRequest := &github.IssueRequest{
		Title:  &title,
		Body:   &description,
		Labels: &labels,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	githubIssue, response, err := client.Issues.Create(ctx, owner, name, issueRequest)
	logging.SharedInstance().MethodInfo("hub", "CreateGithubIssue").Debugf("response of creating github issue: %+v\n", response)
	if err != nil {
		return nil, errors.Wrap(err, "response is error")
	}
	logging.SharedInstance().MethodInfo("hub", "CreateGithubIssue").Debugf("github issue is created: %+v", githubIssue)
	return githubIssue, nil
}

// EditGithubIssue get a issue information from github
func EditGithubIssue(token, owner, name, title, description, state string, issueNumber int, labels []string) (bool, error) {
	client := prepareClient(token)

	issueRequest := &github.IssueRequest{
		Title:  &title,
		Body:   &description,
		State:  &state,
		Labels: &labels,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	issue, response, err := client.Issues.Edit(ctx, owner, name, issueNumber, issueRequest)
	logging.SharedInstance().MethodInfo("hub", "EditGithubIssue").Debugf("response of edit github issue: %+v", response)
	if err != nil {
		return false, errors.Wrap(err, "response is error")
	}
	logging.SharedInstance().MethodInfo("hub", "EditGithubIssue").Debugf("github issue is updated: %+v", issue)
	return true, nil
}

// GetGithubIssue get a issue from github
func GetGithubIssue(token, owner, name string, number int) (*github.Issue, error) {
	client := prepareClient(token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	issue, _, err := client.Issues.Get(ctx, owner, name, number)
	if err != nil {
		return nil, errors.Wrap(err, "response is error")
	}
	return issue, nil
}

// GetGithubIssues returns open and closed issues
func GetGithubIssues(token, owner, name string) ([]*github.Issue, []*github.Issue, error) {
	client := prepareClient(token)

	openIssueOption := github.IssueListByRepoOptions{
		State: "open",
	}
	closedIssueOption := github.IssueListByRepoOptions{
		State: "closed",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opneIssues, _, err := client.Issues.ListByRepo(ctx, owner, name, &openIssueOption)
	if err != nil {
		return nil, nil, errors.Wrap(err, "response is error")
	}
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	closedIssues, _, err := client.Issues.ListByRepo(ctx, owner, name, &closedIssueOption)
	if err != nil {
		return nil, nil, errors.Wrap(err, "response is error")
	}

	return opneIssues, closedIssues, nil
}

// ListLabels list all github labels
func ListLabels(token, owner, name string) ([]*github.Label, error) {
	client := prepareClient(token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	labels, _, err := client.Issues.ListLabels(ctx, owner, name, nil)
	if err != nil {
		return nil, errors.Wrap(err, "response is error")
	}
	return labels, nil
}

// IsPullRequest return true when issue is pull request
func IsPullRequest(issue *github.Issue) bool {
	if issue.PullRequestLinks == nil {
		return false
	}
	return true
}

// CreateWebhook create a new webhook in a github repository
func CreateWebhook(token, owner, name, secret, url string) error {
	client := prepareClient(token)

	hook := webHook(secret, url)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, _, err := client.Repositories.CreateHook(ctx, owner, name, hook)
	if err != nil {
		return errors.Wrap(err, "CreateWebhook error")
	}
	return nil
}

// EditWebhook update a webhook in a github repository
func EditWebhook(token, owner, name, secret, url string, hook *github.Hook) error {
	client := prepareClient(token)

	editHook := webHook(secret, url)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, _, err := client.Repositories.EditHook(ctx, owner, name, *hook.ID, editHook)
	if err != nil {
		return errors.Wrap(err, "EditWebhook error")
	}
	return nil
}

// ListWebhooks list all webhooks in github repository
func ListWebhooks(token, owner, name string) ([]*github.Hook, error) {
	client := prepareClient(token)

	listOptions := &github.ListOptions{
		Page:    1,
		PerPage: 100,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	hooks, _, err := client.Repositories.ListHooks(ctx, owner, name, listOptions)
	if err != nil {
		return nil, errors.Wrap(err, "ListHooks error")
	}
	return hooks, nil
}

// DeleteWebhook delete a webhook in github repository
func DeleteWebhook(token, owner, name string, hook *github.Hook) error {
	client := prepareClient(token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := client.Repositories.DeleteHook(ctx, owner, name, *hook.ID)
	if err != nil {
		return errors.Wrap(err, "DeleteWebhook error")
	}
	return nil
}

func prepareClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	return client
}

func webHook(secret, url string) *github.Hook {
	hookName := "web"
	active := true
	hookConfig := map[string]interface{}{
		"url":          url,
		"content_type": "json",
		"secret":       secret,
	}

	return &github.Hook{
		Name: &hookName,
		URL:  &url,
		Events: []string{
			"commit_comment",
			"push",
			"status",
			"release",
			"issues",
			"issue_comment",
			"pull_request",
			"pull_request_review_comment",
		},
		Active: &active,
		Config: hookConfig,
	}
}

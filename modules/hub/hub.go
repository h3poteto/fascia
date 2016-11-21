package hub

import (
	"github.com/h3poteto/fascia/models/repository"
	"github.com/h3poteto/fascia/modules/logging"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type hub interface {
}

type HubStruct struct {
}

func CheckLabelPresent(token string, repo *repository.RepositoryStruct, title *string) (*github.Label, error) {
	client := prepareClient(token)

	if title == nil {
		return nil, errors.New("title is required")
	}
	githubLabel, response, err := client.Issues.GetLabel(repo.Owner.String, repo.Name.String, *title)
	logging.SharedInstance().MethodInfo("hub", "CheckLabelPresent").Debugf("respone of geting github label: %+v", response)
	if err != nil {
		logging.SharedInstance().MethodInfo("hub", "CheckLabelPresent").Debugf("cannot find github label: %v", repo.Name.String)
		// TODO: 本当はerrorで返したいが，error=nil, label=nilでラベルの存在判定をしている箇所があるので，それらを駆逐したい
		return nil, nil
	}
	logging.SharedInstance().MethodInfo("hub", "CheckLabelPresent").Debugf("github label is exist: %+v", githubLabel)
	return githubLabel, nil
}

func CreateGithubLabel(token string, repo *repository.RepositoryStruct, title *string, color *string) (*github.Label, error) {
	client := prepareClient(token)

	if title == nil || color == nil {
		return nil, errors.New("title or color is nil")
	}

	label := &github.Label{
		Name:  title,
		Color: color,
	}
	githubLabel, response, err := client.Issues.CreateLabel(repo.Owner.String, repo.Name.String, label)
	logging.SharedInstance().MethodInfo("hub", "CreateGithubLabel").Debugf("response of creating github label: %+v\n", response)
	if err != nil {
		return nil, errors.Wrap(err, "response is error")
	}
	logging.SharedInstance().MethodInfo("hub", "CreateGithubLabel").Debugf("github label is created: %+v", githubLabel)
	return githubLabel, nil
}

func UpdateGithubLabel(token string, repo *repository.RepositoryStruct, originalTitle *string, title *string, color *string) (*github.Label, error) {
	client := prepareClient(token)

	if title == nil || color == nil {
		return nil, errors.New("title or color is nil")
	}

	label := &github.Label{
		Name:  title,
		Color: color,
	}
	githubLabel, response, err := client.Issues.EditLabel(repo.Owner.String, repo.Name.String, *originalTitle, label)
	logging.SharedInstance().MethodInfo("hub", "UpddateGithubLabel").Debugf("response of updating github label: %+v", response)
	if err != nil {
		return nil, errors.Wrap(err, "response is error")
	}
	logging.SharedInstance().MethodInfo("hub", "UpdateGithubLabel").Debugf("github label is updated: %+v", githubLabel)
	return githubLabel, nil
}

func CreateGithubIssue(token string, repo *repository.RepositoryStruct, labels []string, title *string, description *string) (*github.Issue, error) {
	client := prepareClient(token)

	if title == nil {
		return nil, errors.New("title is nil")
	}

	issueRequest := &github.IssueRequest{
		Title:  title,
		Body:   description,
		Labels: &labels,
	}
	githubIssue, response, err := client.Issues.Create(repo.Owner.String, repo.Name.String, issueRequest)
	logging.SharedInstance().MethodInfo("hub", "CreateGithubIssue").Debugf("response of creating github issue: %+v\n", response)
	if err != nil {
		return nil, errors.Wrap(err, "response is error")
	}
	logging.SharedInstance().MethodInfo("hub", "CreateGithubIssue").Debugf("github issue is created: %+v", githubIssue)
	return githubIssue, nil
}

// EditGithubIssue get a issue information from github
func EditGithubIssue(token string, repo *repository.RepositoryStruct, issueNumber int, labels []string, title *string, description *string, state *string) (bool, error) {
	client := prepareClient(token)

	issueRequest := &github.IssueRequest{
		Title:  title,
		Body:   description,
		State:  state,
		Labels: &labels,
	}

	issue, response, err := client.Issues.Edit(repo.Owner.String, repo.Name.String, issueNumber, issueRequest)
	logging.SharedInstance().MethodInfo("hub", "EditGithubIssue").Debugf("response of edit github issue: %+v", response)
	if err != nil {
		return false, errors.Wrap(err, "response is error")
	}
	logging.SharedInstance().MethodInfo("hub", "EditGithubIssue").Debugf("github issue is updated: %+v", issue)
	return true, nil
}

// GetGithubIssue get a issue from github
func GetGithubIssue(token string, repo *repository.RepositoryStruct, number int) (*github.Issue, error) {
	client := prepareClient(token)

	issue, _, err := client.Issues.Get(repo.Owner.String, repo.Name.String, number)
	if err != nil {
		return nil, errors.Wrap(err, "response is error")
	}
	return issue, nil
}

func GetGithubIssues(token string, repo *repository.RepositoryStruct) ([]*github.Issue, []*github.Issue, error) {
	client := prepareClient(token)

	openIssueOption := github.IssueListByRepoOptions{
		State: "open",
	}
	closedIssueOption := github.IssueListByRepoOptions{
		State: "closed",
	}
	opneIssues, _, err := client.Issues.ListByRepo(repo.Owner.String, repo.Name.String, &openIssueOption)
	if err != nil {
		return nil, nil, errors.Wrap(err, "response is error")
	}
	closedIssues, _, err := client.Issues.ListByRepo(repo.Owner.String, repo.Name.String, &closedIssueOption)
	if err != nil {
		return nil, nil, errors.Wrap(err, "response is error")
	}

	return opneIssues, closedIssues, nil
}

// ListLabels list all github labels
func ListLabels(token string, repo *repository.RepositoryStruct) ([]*github.Label, error) {
	client := prepareClient(token)

	labels, _, err := client.Issues.ListLabels(repo.Owner.String, repo.Name.String, nil)
	if err != nil {
		return nil, errors.Wrap(err, "response is error")
	}
	return labels, nil
}

func IsPullRequest(issue *github.Issue) bool {
	if issue.PullRequestLinks == nil {
		return false
	}
	return true
}

// CreateWebhook create a new webhook in a github repository
func CreateWebhook(token string, repo *repository.RepositoryStruct, secret string, url string) error {
	client := prepareClient(token)

	name := "web"
	active := true
	hookConfig := map[string]interface{}{
		"url":          url,
		"content_type": "json",
		"secret":       secret,
	}

	hook := github.Hook{
		Name: &name,
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
	_, _, err := client.Repositories.CreateHook(repo.Owner.String, repo.Name.String, &hook)
	if err != nil {
		return errors.Wrap(err, "response is error")
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

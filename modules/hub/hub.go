package hub

import (
	"../../models/repository"
	"../logging"
	"errors"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type hub interface {
}

type HubStruct struct {
}

func CheckLabelPresent(token string, repo *repository.RepositoryStruct, title *string) (*github.Label, error) {
	client := prepareClient(token)

	if title == nil {
		logging.SharedInstance().MethodInfo("hub", "CheckLabelPresent", true).Error("title is nil")
		return nil, errors.New("title is nil")
	}
	githubLabel, response, err := client.Issues.GetLabel(repo.Owner.String, repo.Name.String, *title)
	logging.SharedInstance().MethodInfo("hub", "CheckLabelPresent").Debugf("respone of geting github label: %+v", response)
	if err != nil {
		logging.SharedInstance().MethodInfo("hub", "CheckLabelPresent").Debugf("cannot find github label: %v", repo.Name.String)
		return nil, nil
	}
	logging.SharedInstance().MethodInfo("hub", "CheckLabelPresent").Debugf("github label is exist: %+v", githubLabel)
	return githubLabel, nil
}

func CreateGithubLabel(token string, repo *repository.RepositoryStruct, title *string, color *string) (*github.Label, error) {
	client := prepareClient(token)

	if title == nil || color == nil {
		logging.SharedInstance().MethodInfo("hub", "CreateGithubLabel", true).Error("title or color is nil")
		return nil, errors.New("title or color is nil")
	}

	label := &github.Label{
		Name:  title,
		Color: color,
	}
	githubLabel, response, err := client.Issues.CreateLabel(repo.Owner.String, repo.Name.String, label)
	logging.SharedInstance().MethodInfo("hub", "CreateGithubLabel").Debugf("response of creating github label: %+v\n", response)
	if err != nil {
		return nil, err
	}
	logging.SharedInstance().MethodInfo("hub", "CreateGithubLabel").Debugf("github label is created: %+v", githubLabel)
	return githubLabel, nil
}

func UpdateGithubLabel(token string, repo *repository.RepositoryStruct, originalTitle *string, title *string, color *string) (*github.Label, error) {
	client := prepareClient(token)

	if title == nil || color == nil {
		logging.SharedInstance().MethodInfo("hub", "UpdateGithubLabel", true).Error("title or color is nil")
		return nil, errors.New("title or color is nil")
	}

	label := &github.Label{
		Name:  title,
		Color: color,
	}
	githubLabel, response, err := client.Issues.EditLabel(repo.Owner.String, repo.Name.String, *originalTitle, label)
	logging.SharedInstance().MethodInfo("hub", "UpddateGithubLabel").Debugf("response of updating github label: %+v", response)
	if err != nil {
		return nil, err
	}
	logging.SharedInstance().MethodInfo("hub", "UpdateGithubLabel").Debugf("github label is updated: %+v", githubLabel)
	return githubLabel, nil
}

func CreateGithubIssue(token string, repo *repository.RepositoryStruct, labels []string, title *string, description *string) (*github.Issue, error) {
	client := prepareClient(token)

	if title == nil {
		logging.SharedInstance().MethodInfo("hub", "CreateGithubIssue", true).Error("title is nil")
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
		return nil, err
	}
	logging.SharedInstance().MethodInfo("hub", "CreateGithubIssue").Debugf("github issue is created: %+v", githubIssue)
	return githubIssue, nil
}

func EditGithubIssue(token string, repo *repository.RepositoryStruct, number int64, labels []string, title *string, description *string, state *string) (bool, error) {
	client := prepareClient(token)

	issueRequest := &github.IssueRequest{
		Title:  title,
		Body:   description,
		State:  state,
		Labels: &labels,
	}

	issueNumber := int(number)
	issue, response, err := client.Issues.Edit(repo.Owner.String, repo.Name.String, issueNumber, issueRequest)
	logging.SharedInstance().MethodInfo("hub", "EditGithubIssue").Debugf("response of edit github issue: %+v", response)
	if err != nil {
		return false, err
	}
	logging.SharedInstance().MethodInfo("hub", "EditGithubIssue").Debugf("github issue is updated: %+v", issue)
	return true, nil
}

func GetGithubIssue(token string, repo *repository.RepositoryStruct, number int) (*github.Issue, error) {
	client := prepareClient(token)

	issue, _, err := client.Issues.Get(repo.Owner.String, repo.Name.String, number)
	if err != nil {
		return nil, err
	}
	return issue, nil
}

func GetGithubIssues(token string, repo *repository.RepositoryStruct) ([]github.Issue, []github.Issue, error) {
	client := prepareClient(token)

	openIssueOption := github.IssueListByRepoOptions{
		State: "open",
	}
	closedIssueOption := github.IssueListByRepoOptions{
		State: "closed",
	}
	opneIssues, _, err := client.Issues.ListByRepo(repo.Owner.String, repo.Name.String, &openIssueOption)
	if err != nil {
		return nil, nil, err
	}
	closedIssues, _, err := client.Issues.ListByRepo(repo.Owner.String, repo.Name.String, &closedIssueOption)

	return opneIssues, closedIssues, nil
}

func IsPullRequest(issue *github.Issue) bool {
	if issue.PullRequestLinks == nil {
		return false
	}
	return true
}

func prepareClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	return client
}

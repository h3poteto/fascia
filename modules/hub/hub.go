package hub

import (
	"../../models/repository"
	"../logging"
	"errors"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type hub interface {
}

type HubStruct struct {
}

func CheckLabelPresent(token string, repo *repository.RepositoryStruct, title *string) (*github.Label, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	if title == nil {
		logging.SharedInstance().BaseInfo("hub", "CheckLabelPresent").Error("title is nil")
		return nil, errors.New("title is nil")
	}
	githubLabel, response, err := client.Issues.GetLabel(repo.Owner.String, repo.Name.String, *title)
	logging.SharedInstance().BaseInfo("hub", "CheckLabelPresent").Debugf("respone of geting github label: %+v", response)
	if err != nil {
		logging.SharedInstance().BaseInfo("hub", "CheckLabelPresent").Infof("cannot find github label: %v", repo.Name.String)
		return nil, nil
	}
	logging.SharedInstance().BaseInfo("hub", "CheckLabelPresent").Debugf("github label is exist: %+v", githubLabel)
	return githubLabel, nil
}

func CreateGithubLabel(token string, repo *repository.RepositoryStruct, title *string, color *string) (*github.Label, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	if title == nil || color == nil {
		logging.SharedInstance().BaseInfo("hub", "CreateGithubLabel").Error("title or color is nil")
		return nil, errors.New("title or color is nil")
	}

	label := &github.Label{
		Name:  title,
		Color: color,
	}
	githubLabel, response, err := client.Issues.CreateLabel(repo.Owner.String, repo.Name.String, label)
	logging.SharedInstance().BaseInfo("hub", "CreateGithubLabel").Debugf("response of creating github label: %+v\n", response)
	if err != nil {
		return nil, err
	}
	logging.SharedInstance().BaseInfo("hub", "CreateGithubLabel").Debugf("github label is created: %+v", githubLabel)
	return githubLabel, nil
}

func UpdateGithubLabel(token string, repo *repository.RepositoryStruct, originalTitle *string, title *string, color *string) (*github.Label, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	if title == nil || color == nil {
		logging.SharedInstance().BaseInfo("hub", "UpdateGithubLabel").Error("title or color is nil")
		return nil, errors.New("title or color is nil")
	}

	label := &github.Label{
		Name:  title,
		Color: color,
	}
	githubLabel, response, err := client.Issues.EditLabel(repo.Owner.String, repo.Name.String, *originalTitle, label)
	logging.SharedInstance().BaseInfo("hub", "UpddateGithubLabel").Debugf("response of updating github label: %+v\n", response)
	if err != nil {
		return nil, err
	}
	logging.SharedInstance().BaseInfo("hub", "UpdateGithubLabel").Debugf("github label is updated: %+v", githubLabel)
	return githubLabel, nil
}

func CreateGithubIssue(token string, repo *repository.RepositoryStruct, labels []string, title *string) (*github.Issue, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	if title == nil {
		logging.SharedInstance().BaseInfo("hub", "CreateGithubIssue").Error("title is nil")
		return nil, errors.New("title is nil")
	}

	// TODO: description実装時にはbodyにdescriptionを入れる
	description := ""
	issueRequest := &github.IssueRequest{
		Title:  title,
		Body:   &description,
		Labels: &labels,
	}

	githubIssue, _, err := client.Issues.Create(repo.Owner.String, repo.Name.String, issueRequest)
	if err != nil {
		return nil, err
	}
	logging.SharedInstance().BaseInfo("hub", "CreateGithubIssue").Debugf("github issue is created: %+v", githubIssue)
	return githubIssue, nil
}

func ReplaceLabelsForIssue(token string, repo *repository.RepositoryStruct, number int64, labels []string) (bool, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	issueNumber := int(number)
	_, _, err := client.Issues.ReplaceLabelsForIssue(repo.Owner.String, repo.Name.String, issueNumber, labels)
	if err != nil {
		return false, err
	}
	logging.SharedInstance().BaseInfo("hub", "ReplaceLabelsForIssue").Debugf("label of github issue is replaced: %+v", labels)
	return true, nil
}

func GetGithubIssues(token string, repo *repository.RepositoryStruct) ([]github.Issue, []github.Issue, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

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

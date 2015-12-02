package hub

import (
	"../../models/repository"
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
		return nil, errors.New("title is require")
	}
	githubLabel, response, err := client.Issues.GetLabel(repo.Owner.String, repo.Name.String, *title)
	fmt.Printf("get label for github response: %+v\n", response)
	if err != nil {
		fmt.Printf("cannot find github label: %v\n", repo.Name.String)
		return nil, nil
	}
	fmt.Printf("github label: %+v\n", githubLabel)
	return githubLabel, nil
}

func CreateGithubLabel(token string, repo *repository.RepositoryStruct, title *string, color *string) (*github.Label, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	if title == nil || color == nil {
		return nil, errors.New("title or color is nil")
	}

	label := &github.Label{
		Name:  title,
		Color: color,
	}
	githubLabel, response, err := client.Issues.CreateLabel(repo.Owner.String, repo.Name.String, label)
	fmt.Printf("create label for github response: %+v\n", response)
	if err != nil {
		return nil, err
	}
	fmt.Printf("github label created: %+v\n", githubLabel)
	return githubLabel, nil
}

func UpdateGithubLabel(token string, repo *repository.RepositoryStruct, originalTitle *string, title *string, color *string) (*github.Label, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	if title == nil || color == nil {
		return nil, errors.New("title or color is nil")
	}

	label := &github.Label{
		Name:  title,
		Color: color,
	}
	githubLabel, response, err := client.Issues.EditLabel(repo.Owner.String, repo.Name.String, *originalTitle, label)
	fmt.Printf("update label for github response: %+v\n", response)
	if err != nil {
		return nil, err
	}
	fmt.Printf("github label updated: %+v\n", githubLabel)
	return githubLabel, nil
}

func CreateGithubIssue(token string, repo *repository.RepositoryStruct, labels []string, title *string) (*github.Issue, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	if title == nil {
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
	fmt.Printf("github issue created: %+v\n", githubIssue)
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
	fmt.Printf("github issue replaced labels: %+v\n", labels)
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

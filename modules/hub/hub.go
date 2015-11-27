package hub

import (
	"../../models/repository"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type hub interface {
}

type HubStruct struct {
}

func CheckLabelPresent(token string, repo *repository.RepositoryStruct, title *string) *github.Label {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	if title == nil {
		fmt.Printf("title is nil\n")
		return nil
	}
	githubLabel, response, err := client.Issues.GetLabel(repo.Owner.String, repo.Name.String, *title)
	fmt.Printf("get label for github response: %+v\n", response)
	if err != nil {
		fmt.Printf("cannot find github label: %v\n", repo.Name.String)
		return nil
	}
	fmt.Printf("github label: %+v\n", githubLabel)
	return githubLabel
}

func CreateGithubLabel(token string, repo *repository.RepositoryStruct, title *string, color *string) *github.Label {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	if title == nil || color == nil {
		fmt.Printf("title or color is nil\n")
		return nil
	}

	label := &github.Label{
		Name:  title,
		Color: color,
	}
	githubLabel, response, err := client.Issues.CreateLabel(repo.Owner.String, repo.Name.String, label)
	fmt.Printf("create label for github response: %+v\n", response)
	if err != nil {
		panic(err.Error())
		return nil
	}
	fmt.Printf("github label created: %+v\n", githubLabel)
	return githubLabel
}

func UpdateGithubLabel(token string, repo *repository.RepositoryStruct, originalTitle *string, title *string, color *string) *github.Label {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	if title == nil || color == nil {
		fmt.Printf("title or color is nil\n")
		return nil
	}

	label := &github.Label{
		Name:  title,
		Color: color,
	}
	githubLabel, response, err := client.Issues.EditLabel(repo.Owner.String, repo.Name.String, *originalTitle, label)
	fmt.Printf("update label for github response: %+v\n", response)
	if err != nil {
		panic(err.Error())
		return nil
	}
	fmt.Printf("github label updated: %+v\n", githubLabel)
	return githubLabel
}

func CreateGithubIssue(token string, repo *repository.RepositoryStruct, labels []string, title *string) *github.Issue {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	if title == nil {
		fmt.Printf("title is nil\n")
		return nil
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
		panic(err.Error())
		return nil
	}
	fmt.Printf("github issue created: %+v\n", githubIssue)
	return githubIssue
}

func ReplaceLabelsForIssue(token string, repo *repository.RepositoryStruct, number int64, labels []string) bool {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	issueNumber := int(number)
	_, _, err := client.Issues.ReplaceLabelsForIssue(repo.Owner.String, repo.Name.String, issueNumber, labels)
	if err != nil {
		panic(err.Error())
		return false
	}
	fmt.Printf("github issue replaced labels: %+v\n", labels)
	return true
}

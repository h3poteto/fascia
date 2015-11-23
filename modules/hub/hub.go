package hub

import (
	"../../models/db"
	"../../models/repository"
	"database/sql"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type hub interface {
}

type HubStruct struct {
}

func CheckLabelPresent(listId int64, token string, repo *repository.RepositoryStruct) *github.Label {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	o := &db.Database{}
	var database db.DB = o
	table := database.Init()

	var title sql.NullString
	err := table.QueryRow("select title from list where id = ?;", listId).Scan(&title)
	if err != nil || !title.Valid {
		return nil
	}
	githubLabel, response, err := client.Issues.GetLabel(repo.Owner.String, repo.Name.String, title.String)
	fmt.Printf("get label for github response: %+v\n", response)
	if err != nil {
		fmt.Printf("cannot find github label: %v\n", repo.Name.String)
		return nil
	}
	fmt.Printf("github label: %+v\n", githubLabel)
	return githubLabel
}

func CreateGithubLabel(listId int64, token string, repo *repository.RepositoryStruct) *github.Label {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	o := &db.Database{}
	var database db.DB = o
	table := database.Init()

	var title, color sql.NullString
	err := table.QueryRow("select title, color from list where id = ?;", listId).Scan(&title, &color)
	if err != nil || !title.Valid {
		return nil
	}

	label := &github.Label{
		Name:  &title.String,
		Color: &color.String,
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

func UpdateGithubLabel(listId int64, token string, repo *repository.RepositoryStruct) *github.Label {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	o := &db.Database{}
	var database db.DB = o
	table := database.Init()

	var title, color sql.NullString
	err := table.QueryRow("select title, color from list where id = ?;", listId).Scan(&title, &color)
	if err != nil || !title.Valid {
		return nil
	}

	label := &github.Label{
		Name:  &title.String,
		Color: &color.String,
	}
	githubLabel, response, err := client.Issues.EditLabel(repo.Owner.String, repo.Name.String, title.String, label)
	fmt.Printf("update label for github response: %+v\n", response)
	if err != nil {
		panic(err.Error())
		return nil
	}
	fmt.Printf("github label updated: %+v\n", githubLabel)
	return githubLabel
}

func CreateGithubIssue(taskId int64, token string, repo *repository.RepositoryStruct, labels []string) *github.Issue {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	o := &db.Database{}
	var database db.DB = o
	table := database.Init()

	var title sql.NullString
	err := table.QueryRow("select title from tasks where id = ?;", taskId).Scan(&title)
	if err != nil || !title.Valid {
		return nil
	}

	// TODO: description実装時にはbodyにdescriptionを入れる
	description := ""
	issueRequest := &github.IssueRequest{
		Title:  &title.String,
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

package controllers

import (
	"encoding/json"
	"github.com/google/go-github/github"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
	"net/http"
)

type Github struct {
}

func (u *Github) Repositories(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, result := LoginRequired(r)
	encoder := json.NewEncoder(w)
	if !result {
		http.Error(w, "not logined", 401)
		return
	}
	if !current_user.OauthToken.Valid {
		encoder.Encode(nil)
		return
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: current_user.OauthToken.String},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	nextPage := -1
	var repositories []github.Repository
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
		repos, res, err := client.Repositories.List("", repositoryOption)
		nextPage = res.NextPage
		if err != nil {
			error := JsonError{Error: "repository error"}
			encoder.Encode(error)
			return
		}
		repositories = append(repositories, repos...)

	}
	encoder.Encode(repositories)
}

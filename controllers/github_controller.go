package controllers

import (
	"../modules/logging"
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
	current_user, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("GithubController", "Repositories").Errorf("login error: %v", err.Error())
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
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
			logging.SharedInstance().MethodInfo("GithubController", "Repositories").Errorf("repository error: %v", err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
		repositories = append(repositories, repos...)

	}
	encoder.Encode(repositories)
}

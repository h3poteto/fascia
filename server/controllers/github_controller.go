package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"

	"encoding/json"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
)

type Github struct {
}

func (u *Github) Repositories(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("GithubController", "Repositories", c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	if !currentUser.UserEntity.UserModel.OauthToken.Valid {
		logging.SharedInstance().MethodInfo("GithubController", "Repositories", c).Info("user did not have oauth")
		encoder.Encode(nil)
		return
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: currentUser.UserEntity.UserModel.OauthToken.String},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

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
		repos, res, err := client.Repositories.List("", repositoryOption)
		nextPage = res.NextPage
		if err != nil {
			err := errors.Wrap(err, "repository error")
			logging.SharedInstance().MethodInfoWithStacktrace("GithubController", "Repositories", err, c).Error(err)
			http.Error(w, err.Error(), 500)
			return
		}
		repositories = append(repositories, repos...)

	}
	logging.SharedInstance().MethodInfo("GithubController", "Repositories", c).Info("success to get repositories")
	encoder.Encode(repositories)
}

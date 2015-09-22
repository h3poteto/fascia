package controllers
import (
	"net/http"
	"encoding/json"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
)

type Github struct {
}

func (u *Github)Repositories(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, result := LoginRequired(c, w, r)
	encoder := json.NewEncoder(w)
	if !result {
		error := JsonError{Error: "not logined"}
		encoder.Encode(error)
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
	repositoryOption := &github.RepositoryListOptions{
		Type: "all",
		Sort: "full_name",
		Direction: "asc",
	}
	repos, _, err := client.Repositories.List("", repositoryOption)
	if err != nil {
		error := JsonError{Error: "repository error"}
		encoder.Encode(error)
		return
	}
	encoder.Encode(repos)
}

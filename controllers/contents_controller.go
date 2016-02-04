package controllers

import (
	"../modules/logging"
	"github.com/flosch/pongo2"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
	"net/http"
)

type Contents struct {
}

func (u *Contents) About(c web.C, w http.ResponseWriter, r *http.Request) {
	url := githubOauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	token, err := GenerateCSRFToken(c, w, r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ContentsController", "About", true).Errorf("CSRF error: %v", err)
		http.Error(w, "CSRF error", 500)
		return
	}
	tpl, err := pongo2.DefaultSet.FromFile("about.html.tpl")
	if err != nil {
		logging.SharedInstance().MethodInfo("ContentsController", "About", true).Errorf("template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "About fascia", "oauthURL": url, "token": token}, w)
	return
}

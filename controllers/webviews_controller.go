package controllers

import (
	"../modules/logging"

	"github.com/flosch/pongo2"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
	"net/http"
)

type Webviews struct {
}

func (u *Webviews) SignIn(c web.C, w http.ResponseWriter, r *http.Request) {
	url := githubOauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	token, err := GenerateCSRFToken(c, w, r)
	if err != nil {
		logging.SharedInstance().MethodInfo("WebviewsController", "SignIn", true, c).Errorf("CSRF error: %v", err)
		InternalServerError(w, r)
		return
	}

	tpl, err := pongo2.DefaultSet.FromFile("webviews/sign_in.html.tpl")
	if err != nil {
		logging.SharedInstance().MethodInfo("WebviewsController", "SignIn", true, c).Errorf("template error: %v", err)
		InternalServerError(w, r)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "SignIn", "oauthURL": url, "token": token}, w)
	return
}

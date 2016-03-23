package controllers

import (
	"../models/project"
	"../modules/logging"
	"github.com/flosch/pongo2"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
	"net/http"
	"strconv"
)

type Root struct {
}

func (u *Root) Index(c web.C, w http.ResponseWriter, r *http.Request) {
	currentUser, err := LoginRequired(r)
	// ログインしていない場合はaboutページを見せる
	if err != nil {
		logging.SharedInstance().MethodInfo("RootController", "Index").Infof("login error: %v", err)
		u.About(c, w, r)
		return
	}
	// ログインしている場合はダッシュボードへ
	logging.SharedInstance().MethodInfo("RootController", "Index").Info("login success")

	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if projectID != 0 {
		parentProject, err := project.FindProject(projectID)
		if err != nil || parentProject.UserID != currentUser.ID {
			logging.SharedInstance().MethodInfo("RootController", "Index").Warnf("project not found: %v", err)
			NotFound(w, r)
			return
		}
	}
	tpl, err := pongo2.DefaultSet.FromFile("home.html.tpl")
	if err != nil {
		logging.SharedInstance().MethodInfo("RootController", "Index", true).Errorf("template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "Fascia"}, w)
}

func (u *Root) About(c web.C, w http.ResponseWriter, r *http.Request) {
	url := githubOauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	token, err := GenerateCSRFToken(c, w, r)
	if err != nil {
		logging.SharedInstance().MethodInfo("RootController", "About", true).Errorf("CSRF error: %v", err)
		InternalServerError(w, r)
		return
	}
	tpl, err := pongo2.DefaultSet.FromFile("about.html.tpl")
	if err != nil {
		logging.SharedInstance().MethodInfo("RootController", "About", true).Errorf("template error: %v", err)
		InternalServerError(w, r)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "Fascia", "oauthURL": url, "token": token}, w)
	return
}

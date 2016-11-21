package controllers

import (
	"github.com/h3poteto/fascia/models/project"
	"github.com/h3poteto/fascia/modules/logging"

	"net/http"
	"strconv"

	"github.com/flosch/pongo2"
	"github.com/pkg/errors"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
)

type Root struct {
}

func (u *Root) Index(c web.C, w http.ResponseWriter, r *http.Request) {
	currentUser, err := LoginRequired(r)
	// ログインしていない場合はaboutページを見せる
	if err != nil {
		logging.SharedInstance().MethodInfo("RootController", "Index", c).Infof("login error: %v", err)
		u.About(c, w, r)
		return
	}
	// ログインしている場合はダッシュボードへ
	logging.SharedInstance().MethodInfo("RootController", "Index", c).Info("login success")

	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if projectID != 0 {
		parentProject, err := project.FindProject(projectID)
		if err != nil || parentProject.UserID != currentUser.ID {
			logging.SharedInstance().MethodInfo("RootController", "Index", c).Warnf("project not found: %v", err)
			NotFound(w, r)
			return
		}
	}
	tpl, err := pongo2.DefaultSet.FromFile("home.html.tpl")
	if err != nil {
		err := errors.Wrap(err, "template error")
		logging.SharedInstance().MethodInfoWithStacktrace("RootController", "Index", err, c).Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "Fascia"}, w)
}

func (u *Root) About(c web.C, w http.ResponseWriter, r *http.Request) {
	url := githubOauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	token, err := GenerateCSRFToken(c, w, r)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("RootController", "About", err, c).Error(err)
		InternalServerError(w, r)
		return
	}
	tpl, err := pongo2.DefaultSet.FromFile("about.html.tpl")
	if err != nil {
		err := errors.Wrap(err, "template error")
		logging.SharedInstance().MethodInfoWithStacktrace("RootController", "About", err, c).Error(err)
		InternalServerError(w, r)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "Fascia", "oauthURL": url, "token": token}, w)
	return
}

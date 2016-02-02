package controllers

import (
	"../modules/logging"
	"github.com/flosch/pongo2"
	"github.com/zenazn/goji/web"
	"net/http"
)

type Root struct {
}

func (u *Root) Index(c web.C, w http.ResponseWriter, r *http.Request) {
	_, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("RootController", "Index").Infof("login error: %v", err)
		http.Redirect(w, r, "/sign_in", 302)
		return
	}
	logging.SharedInstance().MethodInfo("RootController", "Index").Info("login success")
	tpl, err := pongo2.DefaultSet.FromFile("home.html.tpl")
	if err != nil {
		logging.SharedInstance().MethodInfo("RootController", "Index", true).Errorf("template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "Fascia"}, w)
}

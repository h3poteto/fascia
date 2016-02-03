package controllers

import (
	"../modules/logging"
	"github.com/flosch/pongo2"
	"github.com/zenazn/goji/web"
	"net/http"
)

type Contents struct {
}

func (u *Contents) About(c web.C, w http.ResponseWriter, r *http.Request) {
	tpl, err := pongo2.DefaultSet.FromFile("about.html.tpl")
	if err != nil {
		logging.SharedInstance().MethodInfo("ContentsController", "About", true).Errorf("template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteWriter(pongo2.Context{"title": "About"}, w)
	return
}

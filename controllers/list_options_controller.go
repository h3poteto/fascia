package controllers

import (
	"../models/list_option"
	"../modules/logging"
	"encoding/json"
	"github.com/zenazn/goji/web"
	"net/http"
)

type ListOptions struct {
}

func (u *ListOptions) Index(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Index").Errorf("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(list_option.ListOptionAll())
}

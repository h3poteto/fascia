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

type ListOptionJSONFormat struct {
	ID     int64
	Action string
}

func (u *ListOptions) Index(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListOptionsController", "Index", false, c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}

	encoder := json.NewEncoder(w)
	jsonOptions := make([]*ListOptionJSONFormat, 0)
	for _, o := range list_option.ListOptionAll() {
		jsonOptions = append(jsonOptions, &ListOptionJSONFormat{ID: o.ID, Action: o.Action})
	}
	encoder.Encode(jsonOptions)
	logging.SharedInstance().MethodInfo("ListOptionsController", "Index", false, c).Info("success to get list options")
}

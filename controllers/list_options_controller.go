package controllers

import (
	"encoding/json"
	"github.com/h3poteto/fascia/models/list_option"
	"github.com/h3poteto/fascia/modules/logging"
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
		logging.SharedInstance().MethodInfo("ListOptionsController", "Index", c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}

	encoder := json.NewEncoder(w)
	jsonOptions := make([]*ListOptionJSONFormat, 0)
	listOptionAll, err := list_option.ListOptionAll()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListOptionsController", "Index", err, c).Error(err)
		http.Error(w, "list options error", 500)
		return
	}
	for _, o := range listOptionAll {
		jsonOptions = append(jsonOptions, &ListOptionJSONFormat{ID: o.ID, Action: o.Action})
	}
	encoder.Encode(jsonOptions)
	logging.SharedInstance().MethodInfo("ListOptionsController", "Index", c).Info("success to get list options")
}

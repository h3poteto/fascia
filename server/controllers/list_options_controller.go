package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/entities/list_option"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/views"
	"github.com/zenazn/goji/web"
)

type ListOptions struct {
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

	listOptionAll, err := handlers.ListOptionAll()
	var optionEntities []*list_option.ListOption
	for _, o := range listOptionAll {
		optionEntities = append(optionEntities, o.ListOptionEntity)
	}
	jsonOptions, err := views.ParseListOptionsJSON(optionEntities)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListOptionsController", "Index", err, c).Error(err)
		http.Error(w, "list options error", 500)
		return
	}

	encoder.Encode(jsonOptions)
	logging.SharedInstance().MethodInfo("ListOptionsController", "Index", c).Info("success to get list options")
}

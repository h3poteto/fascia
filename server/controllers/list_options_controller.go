package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/zenazn/goji/web"
)

type ListOptions struct {
}

// ListOptionJSONFormat defined json format for a list option entity
// TODO: renderで使う型，キャストメソッドも別packageにしたい
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
	listOptionAll, err := handlers.ListOptionAll()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListOptionsController", "Index", err, c).Error(err)
		http.Error(w, "list options error", 500)
		return
	}
	for _, o := range listOptionAll {
		jsonOptions = append(jsonOptions, &ListOptionJSONFormat{
			ID:     o.ListOptionEntity.ListOptionModel.ID,
			Action: o.ListOptionEntity.ListOptionModel.Action,
		})
	}
	encoder.Encode(jsonOptions)
	logging.SharedInstance().MethodInfo("ListOptionsController", "Index", c).Info("success to get list options")
}

package controllers
import (
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
	"github.com/zenazn/goji/web"
	"../models/project"
)

type Lists struct {
}

func (u *Lists)Index(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, result := LoginRequired(c, w, r)
	encoder := json.NewEncoder(w)
	if !result {
		error := JsonError{Error: "not logined"}
		encoder.Encode(error)
		return
	}
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	parentProject := project.FindProject(projectID)
	lists := parentProject.Lists()
	encoder.Encode(lists)
	return
}

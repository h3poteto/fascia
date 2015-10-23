package controllers
import (
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
	"github.com/zenazn/goji/web"
	"github.com/goji/param"
	projectModel "../models/project"
	listModel "../models/list"
)

type Lists struct {
}

type NewListForm struct {
	Title string `param:"title"`
}

func (u *Lists)Index(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, result := LoginRequired(r)
	encoder := json.NewEncoder(w)
	if !result {
		error := JsonError{Error: "not logined"}
		encoder.Encode(error)
		return
	}
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	parentProject := projectModel.FindProject(projectID)
	if parentProject == nil {
		error := JsonError{Error: "project not found"}
		encoder.Encode(error)
		return
	}
	lists := parentProject.Lists()
	encoder.Encode(lists)
	return
}

func (u *Lists)Create(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, result := LoginRequired(r)
	encoder := json.NewEncoder(w)
	if !result {
		error := JsonError{Error: "not logined"}
		encoder.Encode(error)
		return
	}
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	parentProject := projectModel.FindProject(projectID)
	if parentProject == nil {
		error := JsonError{Error: "project not found"}
		encoder.Encode(error)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Wrong Form", 400)
		return
	}
	var newListForm NewListForm
	err = param.Parse(r.PostForm, &newListForm)
	if err != nil {
		http.Error(w, "Wrong parameter", 500)
		return
	}
	fmt.Printf("post new list parameter: %+v\n", newListForm)
	list := listModel.NewList(0, projectID, newListForm.Title)
	if !list.Save() {
		error := JsonError{Error: "save failed"}
		encoder.Encode(error)
		return
	}
	encoder.Encode(*list)
}

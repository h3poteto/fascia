package controllers
import (
	"fmt"
	"net/http"
	"encoding/json"
	"database/sql"
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
	Color string `param:"color"`
}

type EditListForm struct {
	Title string `param:"title"`
	Color string `param:"color"`
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
	for _, l := range lists {
		l.ListTasks = l.Tasks()
	}
	encoder.Encode(lists)
	return
}

func (u *Lists)Create(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, result := LoginRequired(r)
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
	list := listModel.NewList(0, projectID, newListForm.Title, newListForm.Color)

	// github同期処理
	if current_user.OauthToken.Valid {
		token := current_user.OauthToken.String
		repo := parentProject.Repository()
		label := list.CheckLabelPresent(token, repo)
		if label == nil {
			// そもそも既に存在しているなんてことはあまりないのでは
			label = list.CreateGithubLabel(token, repo)
			if label == nil {
				error := JsonError{Error: "failed create github label"}
				encoder.Encode(error)
				return
			}
		}
	}
	if !list.Save() {
		error := JsonError{Error: "save failed"}
		encoder.Encode(error)
		return
	}
	encoder.Encode(*list)
}

func (u *Lists)Update(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, result := LoginRequired(r)
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
	listID, _ := strconv.ParseInt(c.URLParams["list_id"], 10, 64)
	targetList := listModel.FindList(projectID, listID)
	if targetList == nil {
		error := JsonError{Error: "list not found"}
		encoder.Encode(error)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Wrong Form", 400)
		return
	}
	var editListForm EditListForm
	err = param.Parse(r.PostForm, &editListForm)
	if err != nil {
		http.Error(w, "Wrong parameter", 500)
		return
	}
	fmt.Printf("post edit list parameter: %+v\n", editListForm)
	targetList.Title = sql.NullString{String: editListForm.Title, Valid: true}
	targetList.Color = sql.NullString{String: editListForm.Color, Valid: true}

	// TODO: github同期処理
	if current_user.OauthToken.Valid {
	}
	if !targetList.Update() {
		error := JsonError{Error: "save failed"}
		encoder.Encode(error)
		return
	}
	targetList.ListTasks = targetList.Tasks()
	encoder.Encode(*targetList)
}

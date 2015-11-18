package controllers

import (
	listModel "../models/list"
	projectModel "../models/project"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/goji/param"
	"github.com/zenazn/goji/web"
	"net/http"
	"strconv"
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

func (u *Lists) Index(c web.C, w http.ResponseWriter, r *http.Request) {
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
	if parentProject == nil && parentProject.UserId.Int64 != current_user.Id {
		http.Error(w, "project not found", 404)
		return
	}
	lists := parentProject.Lists()
	for _, l := range lists {
		l.ListTasks = l.Tasks()
	}
	encoder.Encode(lists)
	return
}

func (u *Lists) Create(c web.C, w http.ResponseWriter, r *http.Request) {
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
	if parentProject == nil && parentProject.UserId.Int64 != current_user.Id {
		http.Error(w, "project not found", 404)
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
	repo := parentProject.Repository()
	if current_user.OauthToken.Valid && repo != nil {
		token := current_user.OauthToken.String
		label := list.CheckLabelPresent(token, repo)
		if label == nil {
			// そもそも既に存在しているなんてことはあまりないのでは
			label = list.CreateGithubLabel(token, repo)
			if label == nil {
				http.Error(w, "failed create github label", 500)
				return
			}
		} else {
			label = list.UpdateGithubLabel(token, repo)
			if label == nil {
				http.Error(w, "failed update github label", 500)
				return
			}
		}
	}
	if !list.Save() {
		http.Error(w, "failed save", 500)
		return
	}
	encoder.Encode(*list)
}

func (u *Lists) Update(c web.C, w http.ResponseWriter, r *http.Request) {
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
	if parentProject == nil && parentProject.UserId.Int64 != current_user.Id {
		http.Error(w, "project not found", 404)
		return
	}
	listID, _ := strconv.ParseInt(c.URLParams["list_id"], 10, 64)
	targetList := listModel.FindList(projectID, listID)
	if targetList == nil {
		http.Error(w, "list not found", 404)
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

	// github同期処理
	repo := parentProject.Repository()
	if current_user.OauthToken.Valid && repo != nil {
		token := current_user.OauthToken.String
		fmt.Printf("repository: %+v\n", repo)
		label := targetList.CheckLabelPresent(token, repo)
		fmt.Printf("find label: %+v\n", label)
		if label == nil {
			// editの場合はほとんどここには入らない
			label = targetList.CreateGithubLabel(token, repo)
			if label == nil {
				http.Error(w, "failed create github label", 500)
				return
			}
		} else {
			label = targetList.UpdateGithubLabel(token, repo)
			if label == nil {
				http.Error(w, "failed update github label", 500)
				return
			}
		}
	}
	if !targetList.Update() {
		http.Error(w, "save failed", 500)
		return
	}
	targetList.ListTasks = targetList.Tasks()
	encoder.Encode(*targetList)
}

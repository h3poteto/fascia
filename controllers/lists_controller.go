package controllers

import (
	listModel "../models/list"
	projectModel "../models/project"
	taskModel "../models/task"
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

type MoveTaskFrom struct {
	FromListId int64 `param:"from_list_id"`
	ToListId   int64 `param:"to_list_id"`
	TaskId     int64 `param:"task_id"`
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
	repo := parentProject.Repository()
	if current_user.OauthToken.Valid && repo != nil {
		token := current_user.OauthToken.String
		label := list.CheckLabelPresent(token, repo)
		if label == nil {
			// そもそも既に存在しているなんてことはあまりないのでは
			label = list.CreateGithubLabel(token, repo)
			if label == nil {
				error := JsonError{Error: "failed create github label"}
				encoder.Encode(error)
				return
			}
		} else {
			label = list.UpdateGithubLabel(token, repo)
			if label == nil {
				error := JsonError{Error: "failed update github label"}
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
				error := JsonError{Error: "failed create github label"}
				encoder.Encode(error)
				return
			}
		} else {
			label = targetList.UpdateGithubLabel(token, repo)
			if label == nil {
				error := JsonError{Error: "failed update github label"}
				encoder.Encode(error)
				return
			}
		}
	}
	if !targetList.Update() {
		error := JsonError{Error: "save failed"}
		encoder.Encode(error)
		return
	}
	targetList.ListTasks = targetList.Tasks()
	encoder.Encode(*targetList)
}

func (u *Lists) MoveTask(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, result := LoginRequired(r)
	encoder := json.NewEncoder(w)
	if !result {
		http.Error(w, "not logined", 401)
		return
	}

	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	parentProject := projectModel.FindProject(projectID)
	if parentProject == nil && parentProject.UserId.Int64 != current_user.Id {
		http.Error(w, "project not found", 400)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "WrongForm", 400)
		return
	}
	var moveTaskFrom MoveTaskFrom
	err = param.Parse(r.PostForm, &moveTaskFrom)
	if err != nil {
		http.Error(w, "Wrong parameter", 500)
		return
	}
	task := taskModel.FindTask(moveTaskFrom.FromListId, moveTaskFrom.TaskId)
	fmt.Printf("post move taks parameter: %+v\n", moveTaskFrom)
	if !task.ChangeList(moveTaskFrom.ToListId) {
		error := JsonError{Error: "list change failed"}
		encoder.Encode(error)
		return
	}
	allLists := parentProject.Lists()
	for _, l := range allLists {
		l.ListTasks = l.Tasks()
	}
	encoder.Encode(allLists)
	return
}

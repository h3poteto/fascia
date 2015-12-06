package controllers

import (
	listModel "../models/list"
	projectModel "../models/project"
	"../modules/logging"
	"encoding/json"
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
	current_user, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Index").Errorf("login error: %v", err.Error())
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	parentProject := projectModel.FindProject(projectID)
	if parentProject == nil || parentProject.UserId.Int64 != current_user.Id {
		logging.SharedInstance().MethodInfo("ListsController", "Index").Error("project not found")
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
	current_user, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Create").Errorf("login error: %v", err.Error())
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	parentProject := projectModel.FindProject(projectID)
	if parentProject == nil || parentProject.UserId.Int64 != current_user.Id {
		logging.SharedInstance().MethodInfo("ListsController", "Create").Error("project not found")
		http.Error(w, "project not found", 404)
		return
	}

	err = r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Create").Errorf("wrong form: %v", err.Error())
		http.Error(w, "Wrong Form", 400)
		return
	}
	var newListForm NewListForm
	err = param.Parse(r.PostForm, &newListForm)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Create").Errorf("wrong parameter: %v", err.Error())
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ListsController", "Create").Debugf("post new list parameter: %+v", newListForm)
	list := listModel.NewList(0, projectID, current_user.Id, newListForm.Title, newListForm.Color)

	repo := parentProject.Repository()
	if !list.Save(repo, &current_user.OauthToken) {
		logging.SharedInstance().MethodInfo("ListsController", "Create").Error("failed save")
		http.Error(w, "failed save", 500)
		return
	}
	encoder.Encode(*list)
}

func (u *Lists) Update(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Update").Errorf("login error: %v", err.Error())
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	parentProject := projectModel.FindProject(projectID)
	if parentProject == nil || parentProject.UserId.Int64 != current_user.Id {
		logging.SharedInstance().MethodInfo("ListsController", "Update").Error("project not found")
		http.Error(w, "project not found", 404)
		return
	}
	listID, _ := strconv.ParseInt(c.URLParams["list_id"], 10, 64)
	targetList := listModel.FindList(projectID, listID)
	if targetList == nil {
		logging.SharedInstance().MethodInfo("ListsController", "Update").Error("list not found")
		http.Error(w, "list not found", 404)
		return
	}

	err = r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Update").Errorf("wrong form: %v", err.Error())
		http.Error(w, "Wrong Form", 400)
		return
	}
	var editListForm EditListForm
	err = param.Parse(r.PostForm, &editListForm)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Update").Errorf("wrong parameter: %v", err.Error())
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ListsController", "Update").Debugf("post edit list parameter: %+v", editListForm)

	repo := parentProject.Repository()
	if !targetList.Update(repo, &current_user.OauthToken, &editListForm.Title, &editListForm.Color) {
		logging.SharedInstance().MethodInfo("ListsController", "Update").Error("save failed")
		http.Error(w, "save failed", 500)
		return
	}
	targetList.ListTasks = targetList.Tasks()
	encoder.Encode(*targetList)
}

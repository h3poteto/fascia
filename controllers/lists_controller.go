package controllers

import (
	listModel "../models/list"
	projectModel "../models/project"
	"../modules/logging"
	"database/sql"
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
	Title  string `param:"title"`
	Color  string `param:"color"`
	Action string `param:"action"`
}

type ListJSONFormat struct {
	ID           int64
	ProjectID    int64
	UserID       int64
	Title        string
	ListTasks    []*TaskJsonFormat
	Color        string
	ListOptionID int64
}

type AllListJSONFormat struct {
	Lists    []*ListJSONFormat
	NoneList *ListJSONFormat
}

func (u *Lists) Index(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Index").Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Index").Errorf("parse error: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	parentProject, err := projectModel.FindProject(projectID)
	if err != nil || parentProject.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("ListsController", "Index").Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	lists := parentProject.Lists()
	var jsonLists []*ListJSONFormat
	for _, l := range lists {
		jsonLists = append(jsonLists, &ListJSONFormat{ID: l.ID, ProjectID: l.ProjectID, UserID: l.UserID, Title: l.Title.String, ListTasks: TaskFormatToJson(l.Tasks()), Color: l.Color.String, ListOptionID: l.ListOptionID.Int64})
	}
	noneList, err := parentProject.NoneList()
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Index").Error(err)
		http.Error(w, "none list not found", 500)
		return
	}
	jsonNoneList := &ListJSONFormat{ID: noneList.ID, ProjectID: noneList.ProjectID, UserID: noneList.UserID, Title: noneList.Title.String, ListTasks: TaskFormatToJson(noneList.Tasks()), Color: noneList.Color.String, ListOptionID: noneList.ListOptionID.Int64}
	jsonAllLists := AllListJSONFormat{Lists: jsonLists, NoneList: jsonNoneList}
	encoder.Encode(jsonAllLists)
	logging.SharedInstance().MethodInfo("ListsController", "Index").Info("success to get lists")
	return
}

func (u *Lists) Create(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Create").Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Create").Errorf("parse error: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	parentProject, err := projectModel.FindProject(projectID)
	if err != nil || parentProject.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("ListsController", "Create").Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}

	err = r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Create", true).Errorf("wrong form: %v", err)
		http.Error(w, "Wrong Form", 400)
		return
	}
	var newListForm NewListForm
	err = param.Parse(r.PostForm, &newListForm)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Create", true).Errorf("wrong parameter: %v", err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ListsController", "Create").Debugf("post new list parameter: %+v", newListForm)
	list := listModel.NewList(0, projectID, currentUser.ID, newListForm.Title, newListForm.Color, sql.NullInt64{})

	repo, _ := parentProject.Repository()
	if err := list.Save(repo, &currentUser.OauthToken); err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Create", true).Errorf("failed save: %v", err)
		http.Error(w, "failed save", 500)
		return
	}
	jsonList := ListJSONFormat{ID: list.ID, ProjectID: list.ProjectID, UserID: list.UserID, Title: list.Title.String, Color: list.Color.String, ListOptionID: list.ListOptionID.Int64}
	encoder.Encode(jsonList)
	logging.SharedInstance().MethodInfo("ListsController", "Create").Info("success to create list")
	return
}

func (u *Lists) Update(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Update").Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Update").Errorf("parse error: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	parentProject, err := projectModel.FindProject(projectID)
	if err != nil || parentProject.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("ListsController", "Update").Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	listID, err := strconv.ParseInt(c.URLParams["list_id"], 10, 64)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Update").Errorf("parse error: %v", err)
		http.Error(w, "list not found", 404)
		return
	}
	targetList, err := listModel.FindList(projectID, listID)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Update", true).Errorf("list not found: %v", err)
		http.Error(w, "list not found", 404)
		return
	}

	err = r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Update", true).Errorf("wrong form: %v", err)
		http.Error(w, "Wrong Form", 400)
		return
	}
	var editListForm EditListForm
	err = param.Parse(r.PostForm, &editListForm)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Update", true).Errorf("wrong parameter: %v", err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ListsController", "Update").Debugf("post edit list parameter: %+v", editListForm)

	repo, _ := parentProject.Repository()
	if err := targetList.Update(repo, &currentUser.OauthToken, &editListForm.Title, &editListForm.Color, &editListForm.Action); err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Update", true).Errorf("save failed: %v", err)
		http.Error(w, "save failed", 500)
		return
	}
	jsonList := ListJSONFormat{ID: targetList.ID, ProjectID: targetList.ProjectID, UserID: targetList.UserID, Title: targetList.Title.String, ListTasks: TaskFormatToJson(targetList.Tasks()), Color: targetList.Color.String, ListOptionID: targetList.ListOptionID.Int64}
	encoder.Encode(jsonList)
	logging.SharedInstance().MethodInfo("ListsController", "Update").Info("success to update list")
	return
}

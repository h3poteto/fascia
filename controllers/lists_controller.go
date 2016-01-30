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
	Id           int64
	ProjectId    int64
	UserId       int64
	Title        string
	ListTasks    []*TaskJsonFormat
	Color        string
	ListOptionId int64
}

type AllListJSONFormat struct {
	Lists    []*ListJSONFormat
	NoneList *ListJSONFormat
}

func (u *Lists) Index(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Index").Errorf("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	parentProject := projectModel.FindProject(projectID)
	if parentProject == nil || parentProject.UserId != current_user.Id {
		logging.SharedInstance().MethodInfo("ListsController", "Index").Error("project not found")
		http.Error(w, "project not found", 404)
		return
	}
	lists := parentProject.Lists()
	var jsonLists []*ListJSONFormat
	for _, l := range lists {
		jsonLists = append(jsonLists, &ListJSONFormat{Id: l.Id, ProjectId: l.ProjectId, UserId: l.UserId, Title: l.Title.String, ListTasks: TaskFormatToJson(l.Tasks()), Color: l.Color.String, ListOptionId: l.ListOptionId.Int64})
	}
	noneList := parentProject.NoneList()
	jsonNoneList := &ListJSONFormat{Id: noneList.Id, ProjectId: noneList.ProjectId, UserId: noneList.UserId, Title: noneList.Title.String, ListTasks: TaskFormatToJson(noneList.Tasks()), Color: noneList.Color.String, ListOptionId: noneList.ListOptionId.Int64}
	jsonAllLists := AllListJSONFormat{Lists: jsonLists, NoneList: jsonNoneList}
	encoder.Encode(jsonAllLists)
	return
}

func (u *Lists) Create(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Create").Errorf("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	parentProject := projectModel.FindProject(projectID)
	if parentProject == nil || parentProject.UserId != current_user.Id {
		logging.SharedInstance().MethodInfo("ListsController", "Create").Error("project not found")
		http.Error(w, "project not found", 404)
		return
	}

	err = r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Create").Errorf("wrong form: %v", err)
		http.Error(w, "Wrong Form", 400)
		return
	}
	var newListForm NewListForm
	err = param.Parse(r.PostForm, &newListForm)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Create").Errorf("wrong parameter: %v", err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ListsController", "Create").Debugf("post new list parameter: %+v", newListForm)
	list := listModel.NewList(0, projectID, current_user.Id, newListForm.Title, newListForm.Color, sql.NullInt64{})

	repo := parentProject.Repository()
	if !list.Save(repo, &current_user.OauthToken) {
		logging.SharedInstance().MethodInfo("ListsController", "Create").Error("failed save")
		http.Error(w, "failed save", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ListsController", "Create").Info("success to create list")
	jsonList := ListJSONFormat{Id: list.Id, ProjectId: list.ProjectId, UserId: list.UserId, Title: list.Title.String, Color: list.Color.String, ListOptionId: list.ListOptionId.Int64}
	encoder.Encode(jsonList)
}

func (u *Lists) Update(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Update").Errorf("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	parentProject := projectModel.FindProject(projectID)
	if parentProject == nil || parentProject.UserId != current_user.Id {
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
		logging.SharedInstance().MethodInfo("ListsController", "Update").Errorf("wrong form: %v", err)
		http.Error(w, "Wrong Form", 400)
		return
	}
	var editListForm EditListForm
	err = param.Parse(r.PostForm, &editListForm)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Update").Errorf("wrong parameter: %v", err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ListsController", "Update").Debugf("post edit list parameter: %+v", editListForm)

	repo := parentProject.Repository()
	if !targetList.Update(repo, &current_user.OauthToken, &editListForm.Title, &editListForm.Color, &editListForm.Action) {
		logging.SharedInstance().MethodInfo("ListsController", "Update").Error("save failed")
		http.Error(w, "save failed", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ListsController", "Update").Info("success to update list")
	jsonList := ListJSONFormat{Id: targetList.Id, ProjectId: targetList.ProjectId, UserId: targetList.UserId, Title: targetList.Title.String, ListTasks: TaskFormatToJson(targetList.Tasks()), Color: targetList.Color.String, ListOptionId: targetList.ListOptionId.Int64}
	encoder.Encode(jsonList)
}

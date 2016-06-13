package controllers

import (
	listModel "../models/list"
	projectModel "../models/project"
	"../modules/logging"

	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/goji/param"
	"github.com/pkg/errors"
	"github.com/zenazn/goji/web"
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
	ListTasks    []*TaskJSONFormat
	Color        string
	ListOptionID int64
	IsHidden     bool
	IsInitList   bool
}

type AllListJSONFormat struct {
	Lists    []*ListJSONFormat
	NoneList *ListJSONFormat
}

func (u *Lists) Index(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Index", c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Index", err, c).Error(err)
		http.Error(w, "project not found", 404)
		return
	}
	parentProject, err := projectModel.FindProject(projectID)
	if err != nil || parentProject.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("ListsController", "Index", c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	lists, err := parentProject.Lists()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Index", err, c).Error(err)
		http.Error(w, "lists not found", 500)
		return
	}
	jsonLists := ListsFormatToJSON(lists)
	noneList, err := parentProject.NoneList()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Index", err, c).Error(err)
		http.Error(w, "none list not found", 500)
		return
	}
	jsonNoneList := ListFormatToJSON(noneList)
	jsonAllLists := AllListJSONFormat{Lists: jsonLists, NoneList: jsonNoneList}
	encoder.Encode(jsonAllLists)
	logging.SharedInstance().MethodInfo("ListsController", "Index", c).Info("success to get lists")
	return
}

func (u *Lists) Create(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Create", c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Create", err, c).Error(err)
		http.Error(w, "project not found", 404)
		return
	}
	parentProject, err := projectModel.FindProject(projectID)
	if err != nil || parentProject.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("ListsController", "Create", c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}

	err = r.ParseForm()
	if err != nil {
		err := errors.Wrap(err, "wrong form")
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Create", err, c).Error(err)
		http.Error(w, "Wrong Form", 400)
		return
	}
	var newListForm NewListForm
	err = param.Parse(r.PostForm, &newListForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Create", err, c).Error(err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ListsController", "Create", c).Debugf("post new list parameter: %+v", newListForm)
	list := listModel.NewList(0, projectID, currentUser.ID, newListForm.Title, newListForm.Color, sql.NullInt64{}, false)

	repo, _ := parentProject.Repository()
	if err := list.Save(repo, &currentUser.OauthToken); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Create", err, c).Error(err)
		http.Error(w, "failed save", 500)
		return
	}
	jsonList := ListFormatToJSON(list)
	encoder.Encode(jsonList)
	logging.SharedInstance().MethodInfo("ListsController", "Create", c).Info("success to create list")
	return
}

func (u *Lists) Update(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Update", c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Update", err, c).Error(err)
		http.Error(w, "project not found", 404)
		return
	}
	parentProject, err := projectModel.FindProject(projectID)
	if err != nil || parentProject.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("ListsController", "Update", c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	listID, err := strconv.ParseInt(c.URLParams["list_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Update", err, c).Error(err)
		http.Error(w, "list not found", 404)
		return
	}
	targetList, err := listModel.FindList(projectID, listID)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Update", err, c).Error(err)
		http.Error(w, "list not found", 404)
		return
	}

	err = r.ParseForm()
	if err != nil {
		err := errors.Wrap(err, "wrong form")
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Update", err, c).Error(err)
		http.Error(w, "Wrong Form", 400)
		return
	}
	var editListForm EditListForm
	err = param.Parse(r.PostForm, &editListForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Update", err, c).Error(err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ListsController", "Update", c).Debugf("post edit list parameter: %+v", editListForm)

	repo, _ := parentProject.Repository()
	if err := targetList.Update(repo, &currentUser.OauthToken, &editListForm.Title, &editListForm.Color, &editListForm.Action); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Update", err, c).Error(err)
		http.Error(w, "save failed", 500)
		return
	}
	jsonList := ListFormatToJSON(targetList)
	encoder.Encode(jsonList)
	logging.SharedInstance().MethodInfo("ListsController", "Update", c).Info("success to update list")
	return
}

// Hide can hide a list
func (u *Lists) Hide(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Hide", c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Hide", err, c).Error(err)
		http.Error(w, "project not found", 404)
		return
	}
	parentProject, err := projectModel.FindProject(projectID)
	if err != nil || parentProject.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("ListsController", "Hide", c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	listID, err := strconv.ParseInt(c.URLParams["list_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Hide", err, c).Error(err)
		http.Error(w, "list not found", 404)
		return
	}
	targetList, err := listModel.FindList(projectID, listID)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Hide", err, c).Error(err)
		http.Error(w, "list not found", 404)
		return
	}

	if err = targetList.Hide(); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Hide", err, c).Error(err)
		http.Error(w, "hide failed", 500)
		return
	}

	// prepare response
	encoder := json.NewEncoder(w)
	lists, err := parentProject.Lists()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Hide", err, c).Error(err)
		http.Error(w, "lists not found", 500)
		return
	}
	jsonLists := ListsFormatToJSON(lists)
	noneList, err := parentProject.NoneList()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Hide", err, c).Error(err)
		http.Error(w, "none list not found", 500)
		return
	}
	jsonNoneList := ListFormatToJSON(noneList)
	jsonAllLists := AllListJSONFormat{Lists: jsonLists, NoneList: jsonNoneList}
	logging.SharedInstance().MethodInfo("ListsController", "Hide", c).Info("success to hide list")
	encoder.Encode(jsonAllLists)
	return
}

// Display can display a list
func (u *Lists) Display(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Display", c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Display", err, c).Error(err)
		http.Error(w, "project not found", 404)
		return
	}
	parentProject, err := projectModel.FindProject(projectID)
	if err != nil || parentProject.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("ListsController", "Display", c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	listID, err := strconv.ParseInt(c.URLParams["list_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Display", err, c).Error(err)
		http.Error(w, "list not found", 404)
		return
	}
	targetList, err := listModel.FindList(projectID, listID)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Display", err, c).Error(err)
		http.Error(w, "list not found", 404)
		return
	}

	if err = targetList.Display(); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Display", err, c).Error(err)
		http.Error(w, "display failed", 500)
		return
	}

	// prepare response
	encoder := json.NewEncoder(w)
	lists, err := parentProject.Lists()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Display", err, c).Error(err)
		http.Error(w, "lists not found", 500)
		return
	}
	jsonLists := ListsFormatToJSON(lists)
	noneList, err := parentProject.NoneList()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Display", err, c).Error(err)
		http.Error(w, "none list not found", 500)
		return
	}
	jsonNoneList := ListFormatToJSON(noneList)
	jsonAllLists := AllListJSONFormat{Lists: jsonLists, NoneList: jsonNoneList}
	logging.SharedInstance().MethodInfo("ListsController", "Display", c).Info("success to display list")
	encoder.Encode(jsonAllLists)
	return
}

// ListsFormatToJSON convert lists models's array to json
func ListsFormatToJSON(lists []*listModel.ListStruct) []*ListJSONFormat {
	var jsonLists []*ListJSONFormat
	for _, l := range lists {
		jsonLists = append(jsonLists, ListFormatToJSON(l))
	}
	return jsonLists
}

// ListFormatToJSON convert a list model to json
func ListFormatToJSON(list *listModel.ListStruct) *ListJSONFormat {
	tasks, err := list.Tasks()
	if err != nil {
		panic(err)
	}
	return &ListJSONFormat{ID: list.ID, ProjectID: list.ProjectID, UserID: list.UserID, Title: list.Title.String, ListTasks: TaskFormatToJSON(tasks), Color: list.Color.String, ListOptionID: list.ListOptionID.Int64, IsHidden: list.IsHidden, IsInitList: list.IsInitList()}
}

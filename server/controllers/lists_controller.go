package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/validators"
	"github.com/h3poteto/fascia/server/views"

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
	Title    string `param:"title"`
	Color    string `param:"color"`
	OptionID int64  `param:"option_id"`
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
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
		logging.SharedInstance().MethodInfo("ListsController", "Index", c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	lists, err := projectService.ProjectEntity.Lists()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Index", err, c).Error(err)
		http.Error(w, "lists not found", 500)
		return
	}
	noneList, err := projectService.ProjectEntity.NoneList()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Index", err, c).Error(err)
		http.Error(w, "none list not found", 500)
		return
	}

	jsonAllLists, err := views.ParseAllListsJSON(noneList, lists)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Index", err, c).Error(err)
		http.Error(w, "parse error", 500)
		return
	}
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
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
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

	valid, err := validators.ListCreateValidation(newListForm.Title, newListForm.Color)
	if err != nil || !valid {
		logging.SharedInstance().MethodInfo("ListsController", "Create", c).Infof("validation error: %v", err)
		http.Error(w, "validation error", 422)
		return
	}

	list := handlers.NewList(0, projectID, currentUser.UserEntity.UserModel.ID, newListForm.Title, newListForm.Color, sql.NullInt64{}, false)

	if err := list.Save(); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Create", err, c).Error(err)
		http.Error(w, "failed save", 500)
		return
	}
	jsonList, err := views.ParseListJSON(list.ListEntity)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Create", err, c).Error(err)
		http.Error(w, "list format error", 500)
		return
	}
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
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
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
	targetList, err := handlers.FindList(projectID, listID)
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

	valid, err := validators.ListUpdateValidation(
		editListForm.Title,
		editListForm.Color,
		editListForm.OptionID,
	)
	if err != nil || !valid {
		logging.SharedInstance().MethodInfo("ListsController", "Create", c).Infof("validation error: %v", err)
		http.Error(w, "validation error", 422)
		return
	}

	if err := targetList.Update(editListForm.Title, editListForm.Color, editListForm.OptionID); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Update", err, c).Error(err)
		http.Error(w, "save failed", 500)
		return
	}
	jsonList, err := views.ParseListJSON(targetList.ListEntity)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Update", err, c).Error(err)
		http.Error(w, "list format error", 500)
		return
	}
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
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
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
	targetList, err := handlers.FindList(projectID, listID)
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
	lists, err := projectService.ProjectEntity.Lists()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Hide", err, c).Error(err)
		http.Error(w, "lists not found", 500)
		return
	}
	noneList, err := projectService.ProjectEntity.NoneList()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Hide", err, c).Error(err)
		http.Error(w, "none list not found", 500)
		return
	}

	jsonAllLists, err := views.ParseAllListsJSON(noneList, lists)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Hide", err, c).Error(err)
		http.Error(w, "parse error", 500)
		return
	}
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
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
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
	targetList, err := handlers.FindList(projectID, listID)
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
	lists, err := projectService.ProjectEntity.Lists()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Display", err, c).Error(err)
		http.Error(w, "lists not found", 500)
		return
	}
	noneList, err := projectService.ProjectEntity.NoneList()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Display", err, c).Error(err)
		http.Error(w, "none list not found", 500)
		return
	}
	jsonAllLists, err := views.ParseAllListsJSON(noneList, lists)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Display", err, c).Error(err)
		http.Error(w, "parse error", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ListsController", "Display", c).Info("success to display list")
	encoder.Encode(jsonAllLists)
	return
}

package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/validators"
	"github.com/h3poteto/fascia/server/views"

	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
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

func (u *Lists) Index(c echo.Context) error {
	currentUser, err := LoginRequired(c)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Index", c).Infof("login error: %v", err)
		return c.JSON(http.StatusUnauthorized, &JSONError{message: "not logined"})
	}
	projectID, err := strconv.ParseInt(c.Param("project_id"), 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Index", err, c).Error(err)
		return c.JSON(http.StatusNotFound, &JSONError{message: "project not found"})
	}
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
		logging.SharedInstance().MethodInfo("ListsController", "Index", c).Warnf("project not found: %v", err)
		return c.JSON(http.StatusNotFound, &JSONError{message: "project not found"})
	}
	lists, err := projectService.ProjectEntity.Lists()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Index", err, c).Error(err)
		return err
	}
	noneList, err := projectService.ProjectEntity.NoneList()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Index", err, c).Error(err)
		return err
	}

	jsonAllLists, err := views.ParseAllListsJSON(noneList, lists)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Index", err, c).Error(err)
		return err
	}
	logging.SharedInstance().MethodInfo("ListsController", "Index", c).Info("success to get lists")
	return c.JSON(http.StatusOK, jsonAllLists)
}

func (u *Lists) Create(c echo.Context) error {
	currentUser, err := LoginRequired(c)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Create", c).Infof("login error: %v", err)
		return c.JSON(http.StatusUnauthorized, &JSONError{message: "not logined"})
	}

	projectID, err := strconv.ParseInt(c.Param("project_id"), 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Create", err, c).Error(err)
		return c.JSON(http.StatusNotFound, &JSONError{message: "project not found"})
	}
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
		logging.SharedInstance().MethodInfo("ListsController", "Create", c).Warnf("project not found: %v", err)
		return c.JSON(http.StatusNotFound, &JSONError{message: "project not found"})
	}

	var newListForm NewListForm
	err = c.Bind(newListForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Create", err, c).Error(err)
		return err
	}
	logging.SharedInstance().MethodInfo("ListsController", "Create", c).Debugf("post new list parameter: %+v", newListForm)

	valid, err := validators.ListCreateValidation(newListForm.Title, newListForm.Color)
	if err != nil || !valid {
		logging.SharedInstance().MethodInfo("ListsController", "Create", c).Infof("validation error: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, &JSONError{message: "validation error"})
	}

	list := handlers.NewList(0, projectID, currentUser.UserEntity.UserModel.ID, newListForm.Title, newListForm.Color, sql.NullInt64{}, false)

	if err := list.Save(); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Create", err, c).Error(err)
		return err
	}
	jsonList, err := views.ParseListJSON(list.ListEntity)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Create", err, c).Error(err)
		return err
	}
	logging.SharedInstance().MethodInfo("ListsController", "Create", c).Info("success to create list")
	return c.JSON(http.StatusOK, jsonList)
}

func (u *Lists) Update(c echo.Context) error {
	currentUser, err := LoginRequired(c)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Update", c).Infof("login error: %v", err)
		return c.JSON(http.StatusUnauthorized, &JSONError{message: "not logined"})
	}

	projectID, err := strconv.ParseInt(c.Param("project_id"), 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Update", err, c).Error(err)
		return c.JSON(http.StatusNotFound, &JSONError{message: "project not found"})
	}
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
		logging.SharedInstance().MethodInfo("ListsController", "Update", c).Warnf("project not found: %v", err)
		return c.JSON(http.StatusNotFound, &JSONError{message: "project not found"})
	}
	listID, err := strconv.ParseInt(c.Param("list_id"), 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Update", err, c).Error(err)
		return c.JSON(http.StatusNotFound, &JSONError{message: "list not found"})
	}
	targetList, err := handlers.FindList(projectID, listID)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Update", err, c).Error(err)
		return c.JSON(http.StatusNotFound, &JSONError{message: "list not found"})
	}

	var editListForm EditListForm
	err = c.Bind(editListForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Update", err, c).Error(err)
		return err
	}
	logging.SharedInstance().MethodInfo("ListsController", "Update", c).Debugf("post edit list parameter: %+v", editListForm)

	valid, err := validators.ListUpdateValidation(
		editListForm.Title,
		editListForm.Color,
		editListForm.OptionID,
	)
	if err != nil || !valid {
		logging.SharedInstance().MethodInfo("ListsController", "Create", c).Infof("validation error: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, &JSONError{message: "validation error"})
	}

	if err := targetList.Update(editListForm.Title, editListForm.Color, editListForm.OptionID); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Update", err, c).Error(err)
		return err
	}
	jsonList, err := views.ParseListJSON(targetList.ListEntity)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Update", err, c).Error(err)
		return err
	}
	logging.SharedInstance().MethodInfo("ListsController", "Update", c).Info("success to update list")
	return c.JSON(http.StatusOK, jsonList)
}

// Hide can hide a list
func (u *Lists) Hide(c echo.Context) error {
	currentUser, err := LoginRequired(c)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Hide", c).Infof("login error: %v", err)
		return c.JSON(http.StatusUnauthorized, &JSONError{message: "not logined"})
	}
	projectID, err := strconv.ParseInt(c.Param("project_id"), 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Hide", err, c).Error(err)
		return c.JSON(http.StatusNotFound, &JSONError{message: "project not found"})
	}
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
		logging.SharedInstance().MethodInfo("ListsController", "Hide", c).Warnf("project not found: %v", err)
		return c.JSON(http.StatusNotFound, &JSONError{message: "project not found"})
	}
	listID, err := strconv.ParseInt(c.Param("list_id"), 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Hide", err, c).Error(err)
		return c.JSON(http.StatusNotFound, &JSONError{message: "list not found"})
	}
	targetList, err := handlers.FindList(projectID, listID)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Hide", err, c).Error(err)
		return c.JSON(http.StatusNotFound, &JSONError{message: "list not found"})
	}

	if err = targetList.Hide(); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Hide", err, c).Error(err)
		return err
	}

	// prepare response
	lists, err := projectService.ProjectEntity.Lists()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Hide", err, c).Error(err)
		return err
	}
	noneList, err := projectService.ProjectEntity.NoneList()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Hide", err, c).Error(err)
		return err
	}

	jsonAllLists, err := views.ParseAllListsJSON(noneList, lists)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Hide", err, c).Error(err)
		return err
	}
	logging.SharedInstance().MethodInfo("ListsController", "Hide", c).Info("success to hide list")
	return c.JSON(http.StatusOK, jsonAllLists)
}

// Display can display a list
func (u *Lists) Display(c echo.Context) error {
	currentUser, err := LoginRequired(c)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListsController", "Display", c).Infof("login error: %v", err)
		return c.JSON(http.StatusUnauthorized, &JSONError{message: "not logined"})
	}
	projectID, err := strconv.ParseInt(c.Param("project_id"), 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Display", err, c).Error(err)
		return c.JSON(http.StatusNotFound, &JSONError{message: "project not found"})
	}
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
		logging.SharedInstance().MethodInfo("ListsController", "Display", c).Warnf("project not found: %v", err)
		return c.JSON(http.StatusNotFound, &JSONError{message: "project not found"})
	}
	listID, err := strconv.ParseInt(c.Param("list_id"), 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Display", err, c).Error(err)
		return c.JSON(http.StatusNotFound, &JSONError{message: "list not found"})
	}
	targetList, err := handlers.FindList(projectID, listID)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Display", err, c).Error(err)
		return c.JSON(http.StatusNotFound, &JSONError{message: "list not found"})
	}

	if err = targetList.Display(); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Display", err, c).Error(err)
		return err
	}

	// prepare response
	lists, err := projectService.ProjectEntity.Lists()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Display", err, c).Error(err)
		return err
	}
	noneList, err := projectService.ProjectEntity.NoneList()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Display", err, c).Error(err)
		return err
	}
	jsonAllLists, err := views.ParseAllListsJSON(noneList, lists)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListsController", "Display", err, c).Error(err)
		return err
	}
	logging.SharedInstance().MethodInfo("ListsController", "Display", c).Info("success to display list")
	return c.JSON(http.StatusOK, jsonAllLists)
}

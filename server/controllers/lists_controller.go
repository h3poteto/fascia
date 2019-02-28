package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/middlewares"
	"github.com/h3poteto/fascia/server/validators"
	"github.com/h3poteto/fascia/server/views"

	"database/sql"
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// Lists is controller struct for lists
type Lists struct {
}

// NewListForm is struct for new list
type NewListForm struct {
	Title string `json:"title" form:"title"`
	Color string `json:"color" form:"color"`
}

// EditListForm is struct for edit list
type EditListForm struct {
	Title    string `json:"title" form:"title"`
	Color    string `json:"color" form:"color"`
	OptionID int64  `json:"option_id,string" form:"option_id"`
}

// Index returns all lists
func (u *Lists) Index(c echo.Context) error {
	pc, ok := c.(*middlewares.ProjectContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	projectService := pc.ProjectService
	lists, err := projectService.ProjectEntity.Lists()
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	noneList, err := projectService.ProjectEntity.NoneList()
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	jsonAllLists, err := views.ParseAllListsJSON(noneList, lists)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Info("success to get lists")
	return c.JSON(http.StatusOK, jsonAllLists)
}

// Create a new list
func (u *Lists) Create(c echo.Context) error {
	pc, ok := c.(*middlewares.ProjectContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	projectService := pc.ProjectService
	currentUser := pc.CurrentUser

	newListForm := new(NewListForm)
	err := c.Bind(newListForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Debugf("post new list parameter: %+v", newListForm)

	valid, err := validators.ListCreateValidation(newListForm.Title, newListForm.Color)
	if err != nil || !valid {
		logging.SharedInstance().Controller(c).Infof("validation error: %v", err)
		return NewValidationError(err, http.StatusUnprocessableEntity, c)
	}

	list := handlers.NewList(0, projectService.ProjectEntity.ID, currentUser.ID, newListForm.Title, newListForm.Color, sql.NullInt64{}, false)

	if err := list.Save(); err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	jsonList, err := views.ParseListJSON(list.ListEntity)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Info("success to create list")
	return c.JSON(http.StatusOK, jsonList)
}

// Update a list
func (u *Lists) Update(c echo.Context) error {
	lc, ok := c.(*middlewares.ListContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	targetList := lc.ListService

	editListForm := new(EditListForm)
	err := c.Bind(editListForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Debugf("post edit list parameter: %+v", editListForm)

	valid, err := validators.ListUpdateValidation(
		editListForm.Title,
		editListForm.Color,
		editListForm.OptionID,
	)
	if err != nil || !valid {
		logging.SharedInstance().Controller(c).Infof("validation error: %v", err)
		return NewValidationError(err, http.StatusUnprocessableEntity, c)
	}

	if err := targetList.Update(editListForm.Title, editListForm.Color, editListForm.OptionID); err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	jsonList, err := views.ParseListJSON(targetList.ListEntity)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Info("success to update list")
	return c.JSON(http.StatusOK, jsonList)
}

// Hide can hide a list
func (u *Lists) Hide(c echo.Context) error {
	lc, ok := c.(*middlewares.ListContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	targetList := lc.ListService
	projectService := lc.ProjectService

	if err := targetList.Hide(); err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	// prepare response
	lists, err := projectService.ProjectEntity.Lists()
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	noneList, err := projectService.ProjectEntity.NoneList()
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	jsonAllLists, err := views.ParseAllListsJSON(noneList, lists)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Info("success to hide list")
	return c.JSON(http.StatusOK, jsonAllLists)
}

// Display can display a list
func (u *Lists) Display(c echo.Context) error {
	lc, ok := c.(*middlewares.ListContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	projectService := lc.ProjectService
	targetList := lc.ListService

	if err := targetList.Display(); err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	// prepare response
	lists, err := projectService.ProjectEntity.Lists()
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	noneList, err := projectService.ProjectEntity.NoneList()
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	jsonAllLists, err := views.ParseAllListsJSON(noneList, lists)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Info("success to display list")
	return c.JSON(http.StatusOK, jsonAllLists)
}

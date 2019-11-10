package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/middlewares"
	"github.com/h3poteto/fascia/server/usecases/board"
	"github.com/h3poteto/fascia/server/validators"
	"github.com/h3poteto/fascia/server/views"

	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
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

	lists, err := board.ProjectLists(pc.Project)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	noneList, err := board.ProjectNoneList(pc.Project)
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
	project := pc.Project
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

	list, err := board.CreateList(project.ID, currentUser.ID, newListForm.Title, newListForm.Color, sql.NullInt64{}, false)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	jsonList, err := views.ParseListJSON(list)
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
	targetList := lc.List

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

	l, err := board.UpdateList(targetList, editListForm.Title, editListForm.Color, editListForm.OptionID)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	jsonList, err := views.ParseListJSON(l)
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
	targetList := lc.List
	project := lc.Project

	if err := board.HideList(targetList); err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	// prepare response
	lists, err := board.ProjectLists(project)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	noneList, err := board.ProjectNoneList(project)
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
	project := lc.Project
	targetList := lc.List

	if err := board.DisplayList(targetList); err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	// prepare response
	lists, err := board.ProjectLists(project)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	noneList, err := board.ProjectNoneList(project)
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

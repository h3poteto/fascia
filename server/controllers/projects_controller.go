package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/domains/project"
	"github.com/h3poteto/fascia/server/middlewares"
	"github.com/h3poteto/fascia/server/usecases/account"
	"github.com/h3poteto/fascia/server/usecases/board"
	"github.com/h3poteto/fascia/server/validators"
	"github.com/h3poteto/fascia/server/views"

	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// Projects is controller struct for projects
type Projects struct {
}

// NewProjectForm is struct for new project
type NewProjectForm struct {
	Title        string `json:"title" form:"title"`
	Description  string `json:"description" form:"description"`
	RepositoryID int64  `json:"repository_id,string" form:"repository_id"`
}

// EditProjectForm is struct for a project
type EditProjectForm struct {
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
}

// SettingsProjectForm is struct for change settings
type SettingsProjectForm struct {
	ShowIssues       bool `json:"show_issues" form:"show_issues"`
	ShowPullRequests bool `json:"show_pull_requests" form:"show_pull_requests"`
}

// Index returns all projects
func (u *Projects) Index(c echo.Context) error {
	uc, ok := c.(*middlewares.LoginContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	currentUser := uc.CurrentUser
	projects, err := account.UserProjects(currentUser)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	var projectEntities []*project.Project
	for _, p := range projects {
		projectEntities = append(projectEntities, p)
	}
	jsonProjects, err := views.ParseProjectsJSON(projectEntities)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	return c.JSON(http.StatusOK, jsonProjects)
}

// Show return a project detail
func (u *Projects) Show(c echo.Context) error {
	pc, ok := c.(*middlewares.ProjectContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	p := pc.Project
	jsonProject, err := views.ParseProjectJSON(p)
	if err != nil {
		logging.SharedInstance().Controller(c).Error(err)
		return err
	}
	return c.JSON(http.StatusOK, jsonProject)
}

// Create a new project
func (u *Projects) Create(c echo.Context) error {
	uc, ok := c.(*middlewares.LoginContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	currentUser := uc.CurrentUser

	newProjectForm := new(NewProjectForm)
	err := c.Bind(newProjectForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Debugf("post new project parameter: %+v", newProjectForm)

	valid, err := validators.ProjectCreateValidation(
		newProjectForm.Title,
		newProjectForm.Description,
		newProjectForm.RepositoryID,
	)
	if err != nil || !valid {
		logging.SharedInstance().Controller(c).Infof("validation error: %v", err)
		return NewValidationError(err, http.StatusUnprocessableEntity, c)
	}

	p, err := board.CreateProject(
		currentUser.ID,
		newProjectForm.Title,
		newProjectForm.Description,
		newProjectForm.RepositoryID,
		currentUser.OauthToken,
	)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	jsonProject, err := views.ParseProjectJSON(p)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Info("success to create project")
	return c.JSON(http.StatusOK, jsonProject)
}

// Update a project
func (u *Projects) Update(c echo.Context) error {
	pc, ok := c.(*middlewares.ProjectContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	p := pc.Project

	editProjectForm := new(EditProjectForm)
	err := c.Bind(editProjectForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Debugf("post edit project parameter: %+v", editProjectForm)

	valid, err := validators.ProjectUpdateValidation(
		editProjectForm.Title,
		editProjectForm.Description,
	)
	if err != nil || !valid {
		logging.SharedInstance().Controller(c).Infof("validation error: %v", err)
		return NewValidationError(err, http.StatusUnprocessableEntity, c)
	}

	if err := board.UpdateProject(p, editProjectForm.Title, editProjectForm.Description, p.ShowIssues, p.ShowPullRequests); err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Info("success to update project")

	jsonProject, err := views.ParseProjectJSON(p)
	if err != nil {
		logging.SharedInstance().Controller(c).Error(err)
		return err
	}
	return c.JSON(http.StatusOK, jsonProject)
}

// Settings update project settings
func (u *Projects) Settings(c echo.Context) error {
	pc, ok := c.(*middlewares.ProjectContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	p := pc.Project

	settingsProjectForm := new(SettingsProjectForm)
	err := c.Bind(settingsProjectForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Debugf("post edit project parameter: %+v", settingsProjectForm)
	if err := board.UpdateProject(
		p,
		p.Title,
		p.Description,
		settingsProjectForm.ShowIssues,
		settingsProjectForm.ShowPullRequests,
	); err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Info("success to update project")

	jsonProject, err := views.ParseProjectJSON(p)
	if err != nil {
		logging.SharedInstance().Controller(c).Error(err)
		return err
	}
	return c.JSON(http.StatusOK, jsonProject)
}

// FetchGithub import tasks and lists from github
func (u *Projects) FetchGithub(c echo.Context) error {
	pc, ok := c.(*middlewares.ProjectContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	p := pc.Project

	_, err := board.FetchGithub(p)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Errorf("github fetch error: %v", err)
		return err
	}
	lists, err := board.ProjectLists(p)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	noneList, err := board.ProjectNoneList(p)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	jsonAllLists, err := views.ParseAllListsJSON(noneList, lists)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Info("success to fetch github")

	return c.JSON(http.StatusOK, jsonAllLists)
}

// Webhook create a new webhook in github repository
func (u *Projects) Webhook(c echo.Context) error {
	pc, ok := c.(*middlewares.ProjectContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	p := pc.Project

	r, err := board.ProjectRepository(p)
	if err != nil {
		logging.SharedInstance().Controller(c).Error(err)
		return NewJSONError(err, http.StatusNotFound, c)
	}
	token, _ := board.OauthTokenFromProject(p)
	err = r.CreateWebhook(token)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Errorf("failed to create webhook: %v", err)
		return err
	}
	logging.SharedInstance().Controller(c).Info("success to create webhook")
	return c.JSON(http.StatusOK, nil)
}

// Destroy delete a project, all lists and tasks related to a project
func (u *Projects) Destroy(c echo.Context) error {
	pc, ok := c.(*middlewares.ProjectContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	p := pc.Project

	err := board.DeleteProject(p.ID)
	if err != nil {
		logging.SharedInstance().Controller(c).Errorf("project destroy error: %v", err)
		return err
	}
	return c.JSON(http.StatusOK, nil)
}

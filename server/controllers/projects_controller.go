package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/entities/project"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/validators"
	"github.com/h3poteto/fascia/server/views"

	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type Projects struct {
}

type NewProjectForm struct {
	Title        string `param:"title"`
	Description  string `param:"description"`
	RepositoryID int    `param:"repository_id"`
}

type EditProjectForm struct {
	Title       string `param:"title"`
	Description string `param:"description"`
}

type SettingsProjectForm struct {
	ShowIssues       bool `param:"show_issues"`
	ShowPullRequests bool `param:"show_pull_requests"`
}

func (u *Projects) Index(c echo.Context) error {
	currentUser, err := LoginRequired(c)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Index", c).Infof("login error: %v", err)
		return NewJSONError(err, http.StatusUnauthorized, c)
	}
	projects, err := currentUser.Projects()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Index", err, c).Error(err)
		return err
	}

	var projectEntities []*project.Project
	for _, p := range projects {
		projectEntities = append(projectEntities, p.ProjectEntity)
	}
	jsonProjects, err := views.ParseProjectsJSON(projectEntities)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Index", err, c).Error(err)
		return err
	}

	return c.JSON(http.StatusOK, jsonProjects)
}

func (u *Projects) Show(c echo.Context) error {
	currentUser, err := LoginRequired(c)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Show", c).Infof("login error: %v", err)
		return NewJSONError(err, http.StatusUnauthorized, c)
	}
	projectID, err := strconv.ParseInt(c.Param("project_id"), 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Show", err, c).Error(err)
		return NewJSONError(err, http.StatusNotFound, c)
	}
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
		logging.SharedInstance().MethodInfo("ProjectsController", "Show", c).Warnf("project not found: %v", err)
		return NewJSONError(err, http.StatusNotFound, c)
	}

	jsonProject, err := views.ParseProjectJSON(projectService.ProjectEntity)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Show", c).Error(err)
		return err
	}
	return c.JSON(http.StatusOK, jsonProject)
}

func (u *Projects) Create(c echo.Context) error {
	currentUser, err := LoginRequired(c)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Create", c).Infof("login error: %v", err)
		return NewJSONError(err, http.StatusUnauthorized, c)
	}

	var newProjectForm NewProjectForm
	err = c.Bind(newProjectForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Create", err, c).Error(err)
		return err
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Create", c).Debugf("post new project parameter: %+v", newProjectForm)

	valid, err := validators.ProjectCreateValidation(
		newProjectForm.Title,
		newProjectForm.Description,
		newProjectForm.RepositoryID,
	)
	if err != nil || !valid {
		logging.SharedInstance().MethodInfo("ProjectsController", "Create", c).Infof("validation error: %v", err)
		return NewJSONError(err, http.StatusUnprocessableEntity, c)
	}

	projectService, err := handlers.CreateProject(
		currentUser.UserEntity.UserModel.ID,
		newProjectForm.Title,
		newProjectForm.Description,
		newProjectForm.RepositoryID,
		currentUser.UserEntity.UserModel.OauthToken,
	)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Create", err, c).Error(err)
		return err
	}

	jsonProject, err := views.ParseProjectJSON(projectService.ProjectEntity)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Create", err, c).Error(err)
		return err
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Create", c).Info("success to create project")
	return c.JSON(http.StatusOK, jsonProject)
}

// updateはrepository側の更新なしでいこう
// そうしないと，タイトル編集できるって不一致が起こる
func (u *Projects) Update(c echo.Context) error {
	currentUser, err := LoginRequired(c)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update", c).Infof("login error: %v", err)
		return NewJSONError(err, http.StatusUnauthorized, c)
	}
	projectID, err := strconv.ParseInt(c.Param("project_id"), 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Update", err, c).Error(err)
		return NewJSONError(err, http.StatusNotFound, c)
	}
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update", c).Warnf("project not found: %v", err)
		return NewJSONError(err, http.StatusNotFound, c)
	}

	var editProjectForm EditProjectForm
	err = c.Bind(editProjectForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Update", err, c).Error(err)
		return err
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Update", c).Debug("post edit project parameter: %+v", editProjectForm)

	valid, err := validators.ProjectUpdateValidation(
		editProjectForm.Title,
		editProjectForm.Description,
	)
	if err != nil || !valid {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update", c).Infof("validation error: %v", err)
		return NewJSONError(err, http.StatusUnprocessableEntity, c)
	}

	if err := projectService.Update(editProjectForm.Title, editProjectForm.Description, projectService.ProjectEntity.ProjectModel.ShowIssues, projectService.ProjectEntity.ProjectModel.ShowPullRequests); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Update", err, c).Error(err)
		return err
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Update", c).Info("success to update project")

	jsonProject, err := views.ParseProjectJSON(projectService.ProjectEntity)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update", c).Error(err)
		return err
	}
	return c.JSON(http.StatusOK, jsonProject)
}

func (u *Projects) Settings(c echo.Context) error {
	currentUser, err := LoginRequired(c)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Settings", c).Infof("login error: %v", err)
		return NewJSONError(err, http.StatusUnauthorized, c)
	}
	projectID, err := strconv.ParseInt(c.Param("project_id"), 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Settings", err, c).Error(err)
		return NewJSONError(err, http.StatusNotFound, c)
	}
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
		logging.SharedInstance().MethodInfo("ProjectsController", "Settings", c).Warnf("project not found: %v", err)
		return NewJSONError(err, http.StatusNotFound, c)
	}

	var settingsProjectForm SettingsProjectForm
	err = c.Bind(settingsProjectForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Settings", err, c).Error(err)
		return err
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Settings", c).Debug("post edit project parameter: %+v", settingsProjectForm)
	if err := projectService.Update(
		projectService.ProjectEntity.ProjectModel.Title,
		projectService.ProjectEntity.ProjectModel.Description,
		settingsProjectForm.ShowIssues,
		settingsProjectForm.ShowPullRequests,
	); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Settings", err, c).Error(err)
		return err
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Settings", c).Info("success to update project")

	jsonProject, err := views.ParseProjectJSON(projectService.ProjectEntity)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Settings", c).Error(err)
		return err
	}
	return c.JSON(http.StatusOK, jsonProject)
}

func (u *Projects) FetchGithub(c echo.Context) error {
	currentUser, err := LoginRequired(c)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "FetchGithub", c).Infof("login error: %v", err)
		return NewJSONError(err, http.StatusUnauthorized, c)
	}
	projectID, err := strconv.ParseInt(c.Param("project_id"), 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "FetchGithub", err, c).Error(err)
		return NewJSONError(err, http.StatusNotFound, c)
	}
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
		logging.SharedInstance().MethodInfo("ProjectsController", "FetchGithub", c).Warnf("project not found: %v", err)
		return NewJSONError(err, http.StatusNotFound, c)
	}
	_, err = projectService.FetchGithub()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "FetchGithub", err, c).Errorf("github fetch error: %v", err)
		return err
	}
	lists, err := projectService.ProjectEntity.Lists()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "FetchGithub", err, c).Error(err)
		return err
	}
	noneList, err := projectService.ProjectEntity.NoneList()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "FetchGithub", err, c).Error(err)
		return err
	}
	jsonAllLists, err := views.ParseAllListsJSON(noneList, lists)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "FetchGithub", err, c).Error(err)
		return err
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "FetchGithub", c).Info("success to fetch github")

	return c.JSON(http.StatusOK, jsonAllLists)
}

// Webhook create a new webhook in github repository
func (u *Projects) Webhook(c echo.Context) error {
	currentUser, err := LoginRequired(c)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Webhook", c).Infof("login error: %v", err)
		return NewJSONError(err, http.StatusUnauthorized, c)
	}
	projectID, err := strconv.ParseInt(c.Param("project_id"), 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Webhook", err, c).Error(err)
		return NewJSONError(err, http.StatusNotFound, c)
	}
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
		logging.SharedInstance().MethodInfo("ProjectsController", "Webhook", c).Warnf("project not found: %v", err)
		return NewJSONError(err, http.StatusNotFound, c)
	}

	_, find, err := projectService.ProjectEntity.Repository()
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Webhook", c).Error(err)
		return err
	}
	if !find {
		logging.SharedInstance().MethodInfo("ProjectsController", "Webhook", c).Warn("repository not found: %v", err)
		return NewJSONError(err, http.StatusNotFound, c)
	}
	err = projectService.CreateWebhook()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Webhook", err, c).Errorf("failed to create webhook: %v", err)
		return err
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Webhook", c).Info("success to create webhook")
	return c.JSON(http.StatusOK, nil)
}

// Destroy delete a project, all lists and tasks related to a project
func (u *Projects) Destroy(c echo.Context) error {
	currentUser, err := LoginRequired(c)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Destroy", c).Infof("login error: %v", err)
		return NewJSONError(err, http.StatusNotFound, c)
	}
	projectID, err := strconv.ParseInt(c.Param("project_id"), 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Destroy", err, c).Error(err)
		return NewJSONError(err, http.StatusNotFound, c)
	}
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
		logging.SharedInstance().MethodInfo("ProjectsController", "Destroy", c).Warnf("project not found: %v", err)
		return NewJSONError(err, http.StatusNotFound, c)
	}

	err = handlers.DestroyProject(projectID)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Destroy", c).Errorf("project destroy error: %v", err)
		return err
	}
	return c.JSON(http.StatusOK, nil)
}

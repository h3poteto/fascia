package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/entities/project"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/validators"
	"github.com/h3poteto/fascia/server/views"

	"encoding/json"
	"net/http"
	"strconv"

	"github.com/goji/param"
	"github.com/pkg/errors"
	"github.com/zenazn/goji/web"
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

func (u *Projects) Index(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Index", c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projects, err := currentUser.Projects()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Index", err, c).Error(err)
		http.Error(w, "cannot find projects", 500)
		return
	}

	var projectEntities []*project.Project
	for _, p := range projects {
		projectEntities = append(projectEntities, p.ProjectEntity)
	}
	jsonProjects, err := views.ParseProjectsJSON(projectEntities)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Index", err, c).Error(err)
		http.Error(w, "Internal server error", 500)
		return
	}

	encoder.Encode(jsonProjects)
}

func (u *Projects) Show(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Show", c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Show", err, c).Error(err)
		http.Error(w, "project not found", 404)
		return
	}
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
		logging.SharedInstance().MethodInfo("ProjectsController", "Show", c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}

	jsonProject, err := views.ParseProjectJSON(projectService.ProjectEntity)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Show", c).Error(err)
		http.Error(w, "Internal server error", 500)
		return
	}
	encoder.Encode(jsonProject)
	return
}

func (u *Projects) Create(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Create", c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)

	err = r.ParseForm()
	if err != nil {
		err := errors.Wrap(err, "wrong form")
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Create", err, c).Error(err)
		http.Error(w, "Wrong Form", 400)
		return
	}
	var newProjectForm NewProjectForm
	err = param.Parse(r.PostForm, &newProjectForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Create", err, c).Error(err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Create", c).Debugf("post new project parameter: %+v", newProjectForm)

	valid, err := validators.ProjectCreateValidation(
		newProjectForm.Title,
		newProjectForm.Description,
		newProjectForm.RepositoryID,
	)
	if err != nil || !valid {
		logging.SharedInstance().MethodInfo("ProjectsController", "Create", c).Infof("validation error: %v", err)
		http.Error(w, "validation error", 422)
		return
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
		http.Error(w, "save failed", 500)
		return
	}

	jsonProject, err := views.ParseProjectJSON(projectService.ProjectEntity)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Create", err, c).Error(err)
		http.Error(w, "Internal server error", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Create", c).Info("success to create project")
	encoder.Encode(jsonProject)
}

// updateはrepository側の更新なしでいこう
// そうしないと，タイトル編集できるって不一致が起こる
func (u *Projects) Update(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update", c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Update", err, c).Error(err)
		http.Error(w, "project not found", 404)
		return
	}
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update", c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}

	encoder := json.NewEncoder(w)

	err = r.ParseForm()
	if err != nil {
		err := errors.Wrap(err, "wrong form")
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Update", err, c).Error(err)
		http.Error(w, "Wrong Form", 400)
	}
	var editProjectForm EditProjectForm
	err = param.Parse(r.PostForm, &editProjectForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Update", err, c).Error(err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Update", c).Debug("post edit project parameter: %+v", editProjectForm)

	valid, err := validators.ProjectUpdateValidation(
		editProjectForm.Title,
		editProjectForm.Description,
	)
	if err != nil || !valid {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update", c).Infof("validation error: %v", err)
		http.Error(w, "validation error", 422)
		return
	}

	if err := projectService.Update(editProjectForm.Title, editProjectForm.Description, projectService.ProjectEntity.ProjectModel.ShowIssues, projectService.ProjectEntity.ProjectModel.ShowPullRequests); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Update", err, c).Error(err)
		http.Error(w, "update failed", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Update", c).Info("success to update project")

	jsonProject, err := views.ParseProjectJSON(projectService.ProjectEntity)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update", c).Error(err)
		http.Error(w, "Internal server error", 500)
		return
	}
	encoder.Encode(jsonProject)
}

func (u *Projects) Settings(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Settings", c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Settings", err, c).Error(err)
		http.Error(w, "project not found", 404)
		return
	}
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
		logging.SharedInstance().MethodInfo("ProjectsController", "Settings", c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}

	encoder := json.NewEncoder(w)

	err = r.ParseForm()
	if err != nil {
		err := errors.Wrap(err, "wrong form")
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Settings", err, c).Error(err)
		http.Error(w, "Wrong Form", 400)
	}
	var settingsProjectForm SettingsProjectForm
	err = param.Parse(r.PostForm, &settingsProjectForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Settings", err, c).Error(err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Settings", c).Debug("post edit project parameter: %+v", settingsProjectForm)
	if err := projectService.Update(
		projectService.ProjectEntity.ProjectModel.Title,
		projectService.ProjectEntity.ProjectModel.Description,
		settingsProjectForm.ShowIssues,
		settingsProjectForm.ShowPullRequests,
	); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Settings", err, c).Error(err)
		http.Error(w, "update failed", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Settings", c).Info("success to update project")

	jsonProject, err := views.ParseProjectJSON(projectService.ProjectEntity)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Settings", c).Error(err)
		http.Error(w, "Internal server error", 500)
		return
	}
	encoder.Encode(jsonProject)
}

func (u *Projects) FetchGithub(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "FetchGithub", c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "FetchGithub", err, c).Error(err)
		http.Error(w, "project not found", 404)
		return
	}
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
		logging.SharedInstance().MethodInfo("ProjectsController", "FetchGithub", c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	_, err = projectService.FetchGithub()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "FetchGithub", err, c).Errorf("github fetch error: %v", err)
		http.Error(w, err.Error(), 500)
		return
	} else {
		lists, err := projectService.ProjectEntity.Lists()
		if err != nil {
			logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "FetchGithub", err, c).Error(err)
			http.Error(w, "lists not found", 500)
			return
		}
		noneList, err := projectService.ProjectEntity.NoneList()
		if err != nil {
			logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "FetchGithub", err, c).Error(err)
			http.Error(w, "none list not found", 500)
			return
		}
		jsonAllLists, err := views.ParseAllListsJSON(noneList, lists)
		if err != nil {
			logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "FetchGithub", err, c).Error(err)
			http.Error(w, "parse error", 500)
			return
		}
		logging.SharedInstance().MethodInfo("ProjectsController", "FetchGithub", c).Info("success to fetch github")
		encoder.Encode(jsonAllLists)
		return
	}
}

// Webhook create a new webhook in github repository
func (u *Projects) Webhook(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Webhook", c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Webhook", err, c).Error(err)
		http.Error(w, "project not found", 404)
		return
	}
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
		logging.SharedInstance().MethodInfo("ProjectsController", "Webhook", c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}

	_, find, err := projectService.ProjectEntity.Repository()
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Webhook", c).Error(err)
		http.Error(w, "Internal server error", 500)
		return
	}
	if !find {
		logging.SharedInstance().MethodInfo("ProjectsController", "Webhook", c).Warn("repository not found: %v", err)
		http.Error(w, "repository not found", 404)
		return
	}
	err = projectService.CreateWebhook()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Webhook", err, c).Errorf("failed to create webhook: %v", err)
		http.Error(w, "cannot create webhook", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Webhook", c).Info("success to create webhook")
	return
}

// Destroy delete a project, all lists and tasks related to a project
func (u *Projects) Destroy(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Destroy", c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Destroy", err, c).Error(err)
		http.Error(w, "project not found", 404)
		return
	}
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
		logging.SharedInstance().MethodInfo("ProjectsController", "Destroy", c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}

	err = handlers.DestroyProject(projectID)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Destroy", c).Errorf("project destroy error: %v", err)
		http.Error(w, "cannot destroy project", 500)
		return
	}
	return
}

package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/validators"

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

type ProjectJSONFormat struct {
	ID               int64
	UserID           int64
	Title            string
	Description      string
	ShowIssues       bool
	ShowPullRequests bool
	RepositoryID     int64
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
	jsonProjects := make([]*ProjectJSONFormat, 0)
	for _, p := range projects {
		var repositoryID int64
		repo, err := p.ProjectAggregation.Repository()
		if err == nil {
			repositoryID = repo.ID
		}
		jsonProjects = append(jsonProjects, &ProjectJSONFormat{
			ID:               p.ProjectAggregation.ProjectModel.ID,
			UserID:           p.ProjectAggregation.ProjectModel.UserID,
			Title:            p.ProjectAggregation.ProjectModel.Title,
			Description:      p.ProjectAggregation.ProjectModel.Description,
			ShowIssues:       p.ProjectAggregation.ProjectModel.ShowIssues,
			ShowPullRequests: p.ProjectAggregation.ProjectModel.ShowPullRequests,
			RepositoryID:     repositoryID,
		})
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
	if err != nil || !(projectService.CheckOwner(currentUser.UserAggregation.UserModel.ID)) {
		logging.SharedInstance().MethodInfo("ProjectsController", "Show", c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	var repoID int64
	repo, err := projectService.ProjectAggregation.Repository()
	if err == nil {
		repoID = repo.ID
	}
	jsonProject := ProjectJSONFormat{
		ID:               projectService.ProjectAggregation.ProjectModel.ID,
		UserID:           projectService.ProjectAggregation.ProjectModel.UserID,
		Title:            projectService.ProjectAggregation.ProjectModel.Title,
		Description:      projectService.ProjectAggregation.ProjectModel.Description,
		ShowIssues:       projectService.ProjectAggregation.ProjectModel.ShowIssues,
		ShowPullRequests: projectService.ProjectAggregation.ProjectModel.ShowPullRequests,
		RepositoryID:     repoID,
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
		currentUser.UserAggregation.UserModel.ID,
		newProjectForm.Title,
		newProjectForm.Description,
		newProjectForm.RepositoryID,
		currentUser.UserAggregation.UserModel.OauthToken,
	)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Create", err, c).Error(err)
		http.Error(w, "save failed", 500)
		return
	}
	var repositoryID int64
	repo, err := projectService.ProjectAggregation.Repository()
	if err == nil {
		repositoryID = repo.ID
	}
	jsonProject := ProjectJSONFormat{
		ID:               projectService.ProjectAggregation.ProjectModel.ID,
		UserID:           projectService.ProjectAggregation.ProjectModel.UserID,
		Title:            projectService.ProjectAggregation.ProjectModel.Title,
		Description:      projectService.ProjectAggregation.ProjectModel.Description,
		ShowIssues:       projectService.ProjectAggregation.ProjectModel.ShowIssues,
		ShowPullRequests: projectService.ProjectAggregation.ProjectModel.ShowPullRequests,
		RepositoryID:     repositoryID,
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
	if err != nil || !(projectService.CheckOwner(currentUser.UserAggregation.UserModel.ID)) {
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

	if err := projectService.UpdateProject(editProjectForm.Title, editProjectForm.Description, projectService.ProjectAggregation.ProjectModel.ShowIssues, projectService.ProjectAggregation.ProjectModel.ShowPullRequests); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Update", err, c).Error(err)
		http.Error(w, "update failed", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Update", c).Info("success to update project")
	var repositoryID int64
	repo, err := projectService.ProjectAggregation.Repository()
	if err == nil {
		repositoryID = repo.ID
	}
	jsonProject := ProjectJSONFormat{
		ID:               projectService.ProjectAggregation.ProjectModel.ID,
		UserID:           projectService.ProjectAggregation.ProjectModel.UserID,
		Title:            projectService.ProjectAggregation.ProjectModel.Title,
		Description:      projectService.ProjectAggregation.ProjectModel.Description,
		ShowIssues:       projectService.ProjectAggregation.ProjectModel.ShowIssues,
		ShowPullRequests: projectService.ProjectAggregation.ProjectModel.ShowPullRequests,
		RepositoryID:     repositoryID,
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
	if err != nil || !(projectService.CheckOwner(currentUser.UserAggregation.UserModel.ID)) {
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
	if err := projectService.UpdateProject(
		projectService.ProjectAggregation.ProjectModel.Title,
		projectService.ProjectAggregation.ProjectModel.Description,
		settingsProjectForm.ShowIssues,
		settingsProjectForm.ShowPullRequests,
	); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "Settings", err, c).Error(err)
		http.Error(w, "update failed", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Settings", c).Info("success to update project")
	var repositoryID int64
	repo, err := projectService.ProjectAggregation.Repository()
	if err == nil {
		repositoryID = repo.ID
	}
	jsonProject := ProjectJSONFormat{
		ID:               projectService.ProjectAggregation.ProjectModel.ID,
		UserID:           projectService.ProjectAggregation.ProjectModel.UserID,
		Title:            projectService.ProjectAggregation.ProjectModel.Title,
		Description:      projectService.ProjectAggregation.ProjectModel.Description,
		ShowIssues:       projectService.ProjectAggregation.ProjectModel.ShowIssues,
		ShowPullRequests: projectService.ProjectAggregation.ProjectModel.ShowPullRequests,
		RepositoryID:     repositoryID,
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
	if err != nil || !(projectService.CheckOwner(currentUser.UserAggregation.UserModel.ID)) {
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
		lists, err := projectService.ProjectAggregation.Lists()
		if err != nil {
			logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "FetchGithub", err, c).Error(err)
			http.Error(w, "lists not found", 500)
			return
		}
		jsonLists, err := ListsFormatToJSON(lists)
		if err != nil {
			logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "FetchGithub", err, c).Error(err)
			http.Error(w, "lists format error", 500)
			return
		}
		noneList, err := projectService.ProjectAggregation.NoneList()
		if err != nil {
			logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "FetchGithub", err, c).Error(err)
			http.Error(w, "none list not found", 500)
			return
		}
		jsonNoneList, err := ListFormatToJSON(noneList)
		if err != nil {
			logging.SharedInstance().MethodInfoWithStacktrace("ProjectsController", "FetchGithub", err, c).Error(err)
			http.Error(w, "list format error", 500)
			return
		}
		jsonAllLists := AllListJSONFormat{Lists: jsonLists, NoneList: jsonNoneList}
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
	if err != nil || projectService.CheckOwner(currentUser.UserAggregation.UserModel.ID) {
		logging.SharedInstance().MethodInfo("ProjectsController", "Webhook", c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}

	_, err = projectService.ProjectAggregation.Repository()
	if err != nil {
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

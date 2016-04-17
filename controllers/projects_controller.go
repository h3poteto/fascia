package controllers

import (
	projectModel "../models/project"
	"../modules/logging"
	"encoding/json"
	"github.com/goji/param"
	"github.com/zenazn/goji/web"
	"net/http"
	"strconv"
)

type Projects struct {
}

type NewProjectForm struct {
	Title           string `param:"title"`
	Description     string `param:"description"`
	RepositoryID    int64  `param:"repositoryID"`
	RepositoryOwner string `param:"repositoryOwner"`
	RepositoryName  string `param:"repositoryName"`
}

type EditProjectForm struct {
	Title       string `param:"title"`
	Description string `param:"description"`
}

type SettingsProjectForm struct {
	ShowIssues       bool `param:"show_issues"`
	ShowPullRequests bool `param:"show_pull_requests"`
}

type ProjectJsonFormat struct {
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
		logging.SharedInstance().MethodInfo("ProjectsController", "Index", false, c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projects := currentUser.Projects()
	jsonProjects := make([]*ProjectJsonFormat, 0)
	for _, p := range projects {
		var repositoryID int64
		repo, err := p.Repository()
		if err == nil {
			repositoryID = repo.ID
		}
		jsonProjects = append(jsonProjects, &ProjectJsonFormat{ID: p.ID, UserID: p.UserID, Title: p.Title, Description: p.Description, ShowIssues: p.ShowIssues, ShowPullRequests: p.ShowPullRequests, RepositoryID: repositoryID})
	}
	encoder.Encode(jsonProjects)
}

func (u *Projects) Show(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Show", false, c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Show", true, c).Errorf("parse error: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	project, err := projectModel.FindProject(projectID)
	if err != nil || project.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("ProjectsController", "Show", false, c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	var repoID int64
	repo, err := project.Repository()
	if err == nil {
		repoID = repo.ID
	}
	jsonProject := ProjectJsonFormat{ID: project.ID, UserID: project.UserID, Title: project.Title, Description: project.Description, ShowIssues: project.ShowIssues, ShowPullRequests: project.ShowPullRequests, RepositoryID: repoID}
	encoder.Encode(jsonProject)
	return
}

func (u *Projects) Create(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Create", false, c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)

	err = r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Create", true, c).Errorf("wrong form: %v", err)
		http.Error(w, "Wrong Form", 400)
		return
	}
	var newProjectForm NewProjectForm
	err = param.Parse(r.PostForm, &newProjectForm)

	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Create", true, c).Errorf("wrong paramter: %v", err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Create", false, c).Debugf("post new project parameter: %+v", newProjectForm)

	project, err := projectModel.Create(currentUser.ID, newProjectForm.Title, newProjectForm.Description, newProjectForm.RepositoryID, newProjectForm.RepositoryOwner, newProjectForm.RepositoryName, currentUser.OauthToken)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Create", true, c).Error("save failed")
		http.Error(w, "save failed", 500)
		return
	}
	var repositoryID int64
	repo, err := project.Repository()
	if err == nil {
		repositoryID = repo.ID
	}
	jsonProject := ProjectJsonFormat{ID: project.ID, UserID: project.UserID, Title: project.Title, Description: project.Description, ShowIssues: project.ShowIssues, ShowPullRequests: project.ShowPullRequests, RepositoryID: repositoryID}
	logging.SharedInstance().MethodInfo("ProjectsController", "Create", false, c).Info("success to create project")
	encoder.Encode(jsonProject)
}

// updateはrepository側の更新なしでいこう
// そうしないと，タイトル編集できるって不一致が起こる
func (u *Projects) Update(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update", false, c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update", true, c).Errorf("parse error: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	project, err := projectModel.FindProject(projectID)
	if err != nil || project.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update", false, c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}

	encoder := json.NewEncoder(w)

	err = r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update", true, c).Errorf("wrong form: %v", err)
		http.Error(w, "Wrong Form", 400)
	}
	var editProjectForm EditProjectForm
	err = param.Parse(r.PostForm, &editProjectForm)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update", true, c).Errorf("wrong parameter: %v", err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Update", false, c).Debug("post edit project parameter: %+v", editProjectForm)
	if err := project.Update(editProjectForm.Title, editProjectForm.Description, project.ShowIssues, project.ShowPullRequests); err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update", true, c).Errorf("update failed: %v", err)
		http.Error(w, "update failed", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Update", false, c).Info("success to update project")
	var repositoryID int64
	repo, err := project.Repository()
	if err == nil {
		repositoryID = repo.ID
	}
	jsonProject := ProjectJsonFormat{ID: project.ID, UserID: project.UserID, Title: project.Title, Description: project.Description, ShowIssues: project.ShowIssues, ShowPullRequests: project.ShowPullRequests, RepositoryID: repositoryID}
	encoder.Encode(jsonProject)
}

func (u *Projects) Settings(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Settings", false, c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Settings", true, c).Errorf("parse error: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	project, err := projectModel.FindProject(projectID)
	if err != nil || project.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("ProjectsController", "Settings", false, c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}

	encoder := json.NewEncoder(w)

	err = r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Settings", true, c).Errorf("wrong form: %v", err)
		http.Error(w, "Wrong Form", 400)
	}
	var settingsProjectForm SettingsProjectForm
	err = param.Parse(r.PostForm, &settingsProjectForm)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Settings", true, c).Errorf("wrong parameter: %v", err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Settings", false, c).Debug("post edit project parameter: %+v", settingsProjectForm)
	if err := project.Update(project.Title, project.Description, settingsProjectForm.ShowIssues, settingsProjectForm.ShowPullRequests); err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Settings", true, c).Errorf("update failed: %v", err)
		http.Error(w, "update failed", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Settings", false, c).Info("success to update project")
	var repositoryID int64
	repo, err := project.Repository()
	if err == nil {
		repositoryID = repo.ID
	}
	jsonProject := ProjectJsonFormat{ID: project.ID, UserID: project.UserID, Title: project.Title, Description: project.Description, ShowIssues: project.ShowIssues, ShowPullRequests: project.ShowPullRequests, RepositoryID: repositoryID}
	encoder.Encode(jsonProject)
}

func (u *Projects) FetchGithub(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "FetchGithub", false, c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "FetchGithub", true, c).Errorf("parse error: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	project, err := projectModel.FindProject(projectID)
	if err != nil || project.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("ProjectsController", "FetchGithub", false, c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	_, err = project.FetchGithub()
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "FetchGithub", true, c).Errorf("github fetch error: %v", err)
		http.Error(w, err.Error(), 500)
		return
	} else {
		lists := project.Lists()
		jsonLists := ListsFormatToJson(lists)
		noneList, err := project.NoneList()
		if err != nil {
			logging.SharedInstance().MethodInfo("ProjectsController", "FetchGithub", true, c).Error(err)
			http.Error(w, "none list not found", 500)
			return
		}
		jsonNoneList := ListFormatToJson(noneList)
		jsonAllLists := AllListJSONFormat{Lists: jsonLists, NoneList: jsonNoneList}
		logging.SharedInstance().MethodInfo("ProjectsController", "FetchGithub", false, c).Info("success to fetch github")
		encoder.Encode(jsonAllLists)
		return
	}
}

// Webhook create a new webhook in github repository
func (u *Projects) Webhook(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Webhook", false, c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Webhook", true, c).Errorf("parse error: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	project, err := projectModel.FindProject(projectID)
	if err != nil || project.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("ProjectsController", "Webhook", false, c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}

	_, err = project.Repository()
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Webhook", false, c).Warn("repository not found: %v", err)
		http.Error(w, "repository not found", 404)
		return
	}
	err = project.CreateWebhook()
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Webhook", true, c).Errorf("failed to create webhook: %v", err)
		http.Error(w, "cannot create webhook", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Webhook", false, c).Info("success to create webhook")
	return
}

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
		logging.SharedInstance().MethodInfo("ProjectsController", "Index").Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projects := currentUser.Projects()
	jsonProjects := make([]*ProjectJsonFormat, 0)
	for _, p := range projects {
		var repositoryID int64
		repo := p.Repository()
		if repo != nil {
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
		logging.SharedInstance().MethodInfo("ProjectsController", "Show").Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	project := projectModel.FindProject(projectID)
	if project == nil || project.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("ProjectsController", "Show").Warn("project not found")
		http.Error(w, "project not found", 404)
		return
	}
	var repoID int64
	repo := project.Repository()
	if repo != nil {
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
		logging.SharedInstance().MethodInfo("ProjectsController", "Create").Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)

	err = r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Create", true).Errorf("wrong form: %v", err)
		http.Error(w, "Wrong Form", 400)
		return
	}
	var newProjectForm NewProjectForm
	err = param.Parse(r.PostForm, &newProjectForm)

	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Create", true).Errorf("wrong paramter: %v", err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Create").Debugf("post new project parameter: %+v", newProjectForm)

	project, err := projectModel.Create(currentUser.ID, newProjectForm.Title, newProjectForm.Description, newProjectForm.RepositoryID, newProjectForm.RepositoryOwner, newProjectForm.RepositoryName, currentUser.OauthToken)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Create", true).Error("save failed")
		http.Error(w, "save failed", 500)
		return
	}
	var repositoryID int64
	repo := project.Repository()
	if repo != nil {
		repositoryID = repo.ID
	}
	jsonProject := ProjectJsonFormat{ID: project.ID, UserID: project.UserID, Title: project.Title, Description: project.Description, ShowIssues: project.ShowIssues, ShowPullRequests: project.ShowPullRequests, RepositoryID: repositoryID}
	logging.SharedInstance().MethodInfo("ProjectsController", "Create").Info("success to create project")
	encoder.Encode(jsonProject)
}

// updateはrepository側の更新なしでいこう
// そうしないと，タイトル編集できるって不一致が起こる
func (u *Projects) Update(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update").Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	project := projectModel.FindProject(projectID)
	if project == nil || project.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update").Warn("project not found")
		http.Error(w, "project not found", 404)
		return
	}

	encoder := json.NewEncoder(w)

	err = r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update", true).Errorf("wrong form: %v", err)
		http.Error(w, "Wrong Form", 400)
	}
	var editProjectForm EditProjectForm
	err = param.Parse(r.PostForm, &editProjectForm)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update", true).Errorf("wrong parameter: %v", err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Update").Debug("post edit project parameter: %+v", editProjectForm)
	if !project.Update(editProjectForm.Title, editProjectForm.Description, project.ShowIssues, project.ShowPullRequests) {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update", true).Error("update failed")
		http.Error(w, "update failed", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Update").Info("success to update project")
	var repositoryID int64
	repo := project.Repository()
	if repo != nil {
		repositoryID = repo.ID
	}
	jsonProject := ProjectJsonFormat{ID: project.ID, UserID: project.UserID, Title: project.Title, Description: project.Description, ShowIssues: project.ShowIssues, ShowPullRequests: project.ShowPullRequests, RepositoryID: repositoryID}
	encoder.Encode(jsonProject)
}

func (u *Projects) Settings(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Settings").Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	project := projectModel.FindProject(projectID)
	if project == nil || project.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("ProjectsController", "Settings").Warn("project not found")
		http.Error(w, "project not found", 404)
		return
	}

	encoder := json.NewEncoder(w)

	err = r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Settings", true).Errorf("wrong form: %v", err)
		http.Error(w, "Wrong Form", 400)
	}
	var settingsProjectForm SettingsProjectForm
	err = param.Parse(r.PostForm, &settingsProjectForm)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Settings", true).Errorf("wrong parameter: %v", err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Settings").Debug("post edit project parameter: %+v", settingsProjectForm)
	if !project.Update(project.Title, project.Description, settingsProjectForm.ShowIssues, settingsProjectForm.ShowPullRequests) {
		logging.SharedInstance().MethodInfo("ProjectsController", "Settings", true).Error("update failed")
		http.Error(w, "update failed", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Settings").Info("success to update project")
	var repositoryID int64
	repo := project.Repository()
	if repo != nil {
		repositoryID = repo.ID
	}
	jsonProject := ProjectJsonFormat{ID: project.ID, UserID: project.UserID, Title: project.Title, Description: project.Description, ShowIssues: project.ShowIssues, ShowPullRequests: project.ShowPullRequests, RepositoryID: repositoryID}
	encoder.Encode(jsonProject)
}

func (u *Projects) FetchGithub(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "FetchGithub").Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	project := projectModel.FindProject(projectID)
	if project == nil || project.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("ProjectsController", "FetchGithub").Warn("project not found")
		http.Error(w, "project not found", 404)
		return
	}
	_, err = project.FetchGithub()
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "FetchGithub", true).Errorf("github fetch error: %v", err)
		http.Error(w, err.Error(), 500)
		return
	} else {
		lists := project.Lists()
		var jsonLists []*ListJSONFormat
		for _, l := range lists {
			jsonLists = append(jsonLists, &ListJSONFormat{ID: l.ID, ProjectID: l.ProjectID, UserID: l.UserID, Title: l.Title.String, ListTasks: TaskFormatToJson(l.Tasks()), Color: l.Color.String, ListOptionID: l.ListOptionID.Int64})
		}
		noneList := project.NoneList()
		jsonNoneList := &ListJSONFormat{ID: noneList.ID, ProjectID: noneList.ProjectID, UserID: noneList.UserID, Title: noneList.Title.String, ListTasks: TaskFormatToJson(noneList.Tasks()), Color: noneList.Color.String, ListOptionID: noneList.ListOptionID.Int64}
		jsonAllLists := AllListJSONFormat{Lists: jsonLists, NoneList: jsonNoneList}
		logging.SharedInstance().MethodInfo("ProjectsController", "FetchGithub").Info("success to fetch github")
		encoder.Encode(jsonAllLists)
		return
	}
}

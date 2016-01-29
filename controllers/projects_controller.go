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
	RepositoryID    int64  `param:"repositoryId"`
	RepositoryOwner string `param:"repositoryOwner"`
	RepositoryName  string `param:"repositoryName"`
}

type EditProjectForm struct {
	Title       string `param:"title"`
	Description string `param:"description"`
}

type ProjectJsonFormat struct {
	Id           int64
	UserId       int64
	Title        string
	Description  string
	RepositoryID int64
}

func (u *Projects) Index(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Index").Errorf("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projects := current_user.Projects()
	jsonProjects := make([]*ProjectJsonFormat, 0)
	for _, p := range projects {
		jsonProjects = append(jsonProjects, &ProjectJsonFormat{Id: p.Id, UserId: p.UserId, Title: p.Title, Description: p.Description})
	}
	encoder.Encode(jsonProjects)
}

func (u *Projects) Show(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Show").Errorf("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	project := projectModel.FindProject(projectID)
	if project == nil || project.UserId != current_user.Id {
		logging.SharedInstance().MethodInfo("ProjectsController", "Show").Error("project not found")
		http.Error(w, "project not found", 404)
		return
	}
	var repoId int64
	repo := project.Repository()
	if repo != nil {
		repoId = repo.Id
	}
	jsonProject := ProjectJsonFormat{Id: project.Id, UserId: project.UserId, Title: project.Title, Description: project.Description, RepositoryID: repoId}
	encoder.Encode(jsonProject)
	return
}

func (u *Projects) Create(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Create").Errorf("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)

	err = r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Create").Errorf("wrong form: %v", err)
		http.Error(w, "Wrong Form", 400)
		return
	}
	var newProjectForm NewProjectForm
	err = param.Parse(r.PostForm, &newProjectForm)

	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Create").Errorf("wrong paramter: %v", err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Create").Debugf("post new project parameter: %+v", newProjectForm)

	project, err := projectModel.Create(current_user.Id, newProjectForm.Title, newProjectForm.Description, newProjectForm.RepositoryID, newProjectForm.RepositoryOwner, newProjectForm.RepositoryName, current_user.OauthToken)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Create").Error("save failed")
		http.Error(w, "save failed", 500)
		return
	}
	jsonProject := ProjectJsonFormat{Id: project.Id, UserId: project.UserId, Title: project.Title, Description: project.Description}
	logging.SharedInstance().MethodInfo("ProjectsController", "Create").Info("success to create project")
	encoder.Encode(jsonProject)
}

// updateはrepository側の更新なしでいこう
// そうしないと，タイトル編集できるって不一致が起こる
func (u *Projects) Update(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update").Errorf("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	project := projectModel.FindProject(projectID)
	if project == nil || project.UserId != current_user.Id {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update").Error("project not found")
		http.Error(w, "project not found", 404)
		return
	}

	encoder := json.NewEncoder(w)

	err = r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update").Errorf("wrong form: %v", err)
		http.Error(w, "Wrong Form", 400)
	}
	var editProjectForm EditProjectForm
	err = param.Parse(r.PostForm, &editProjectForm)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update").Errorf("wrong parameter: %v", err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Update").Debug("post edit project parameter: %+v", editProjectForm)
	if !project.Update(editProjectForm.Title, editProjectForm.Description) {
		logging.SharedInstance().MethodInfo("ProjectsController", "Update").Error("update failed")
		http.Error(w, "update failed", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Update").Info("success to update project")
	var repositoryId int64
	repo := project.Repository()
	if repo != nil {
		repositoryId = repo.Id
	}
	jsonProject := ProjectJsonFormat{Id: project.Id, UserId: project.UserId, Title: project.Title, Description: project.Description, RepositoryID: repositoryId}
	encoder.Encode(jsonProject)
}

func (u *Projects) FetchGithub(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "FetchGithub").Errorf("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	project := projectModel.FindProject(projectID)
	if project == nil || project.UserId != current_user.Id {
		logging.SharedInstance().MethodInfo("ProjectsController", "FetchGithub").Error("project not found")
		http.Error(w, "project not found", 404)
		return
	}
	_, err = project.FetchGithub()
	if err != nil {
		logging.SharedInstance().MethodInfo("ProjectsController", "FetchGithub").Errorf("github fetch error: %v", err)
		http.Error(w, err.Error(), 500)
		return
	} else {
		lists := project.Lists()
		jsonLists := make([]*ListJsonFormat, 0)
		for _, l := range lists {
			jsonLists = append(jsonLists, &ListJsonFormat{Id: l.Id, ProjectId: l.ProjectId, UserId: l.UserId, Title: l.Title.String, ListTasks: TaskFormatToJson(l.Tasks()), Color: l.Color.String, ListOptionId: l.ListOptionId.Int64})
		}
		noneList := project.NoneList()
		jsonNoneList := &ListJsonFormat{Id: noneList.Id, ProjectId: noneList.ProjectId, UserId: noneList.UserId, Title: noneList.Title.String, ListTasks: TaskFormatToJson(noneList.Tasks()), Color: noneList.Color.String, ListOptionId: noneList.ListOptionId.Int64}
		jsonAllLists := AllListJsonFormat{Lists: jsonLists, NoneList: jsonNoneList}
		logging.SharedInstance().MethodInfo("ProjectsController", "FetchGithub").Info("success to fetch github")
		encoder.Encode(jsonAllLists)
		return
	}
}

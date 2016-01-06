package controllers

import (
	listModel "../models/list"
	"../models/list_option"
	projectModel "../models/project"
	repositoryModel "../models/repository"
	"../modules/logging"
	"database/sql"
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

type ProjectJsonFormat struct {
	Id          int64
	UserId      int64
	Title       string
	Description string
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
	jsonProject := ProjectJsonFormat{Id: project.Id, UserId: project.UserId, Title: project.Title, Description: project.Description}
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
	project := projectModel.NewProject(0, current_user.Id, newProjectForm.Title, newProjectForm.Description)
	if !project.Save() {
		logging.SharedInstance().MethodInfo("ProjectsController", "Create").Error("save failed")
		http.Error(w, "save failed", 500)
		return
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Create").Info("success to create project")

	// closeのlist_optionだけはdoneにつけておきたい
	closeListOption := list_option.FindByAction("close")

	// 初期リストを作っておく
	// TODO: ここエラーハンドリングしたほうがいい
	todo := listModel.NewList(0, project.Id, current_user.Id, "ToDo", "ff0000", sql.NullInt64{})
	inprogress := listModel.NewList(0, project.Id, current_user.Id, "InProgress", "0000ff", sql.NullInt64{})
	done := listModel.NewList(0, project.Id, current_user.Id, "Done", "0a0a0a", sql.NullInt64{Int64: closeListOption.Id, Valid: true})
	if newProjectForm.RepositoryID != 0 {
		repository := repositoryModel.NewRepository(0, project.Id, newProjectForm.RepositoryID, newProjectForm.RepositoryOwner, newProjectForm.RepositoryName)
		repository.Save()
		todo.Save(repository, &current_user.OauthToken)
		inprogress.Save(repository, &current_user.OauthToken)
		done.Save(repository, &current_user.OauthToken)
	} else {
		todo.Save(nil, nil)
		inprogress.Save(nil, nil)
		done.Save(nil, nil)
	}
	logging.SharedInstance().MethodInfo("ProjectsController", "Create").Info("success to create initial lists")

	jsonProject := ProjectJsonFormat{Id: project.Id, UserId: project.UserId, Title: project.Title, Description: project.Description}
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
			jsonLists = append(jsonLists, &ListJsonFormat{Id: l.Id, ProjectId: l.ProjectId, UserId: l.UserId, Title: l.Title.String, ListTasks: TaskFormatToJson(l.Tasks()), Color: l.Color.String})
		}
		logging.SharedInstance().MethodInfo("ProjectsController", "FetchGithub").Info("success to fetch github")
		encoder.Encode(jsonLists)
		return
	}
}

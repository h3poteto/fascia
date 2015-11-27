package controllers

import (
	listModel "../models/list"
	projectModel "../models/project"
	repositoryModel "../models/repository"
	"encoding/json"
	"fmt"
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

func (u *Projects) Index(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, result := LoginRequired(r)
	encoder := json.NewEncoder(w)
	if !result {
		http.Error(w, "not logined", 401)
		return
	}
	projects := current_user.Projects()
	encoder.Encode(projects)
}

func (u *Projects) Show(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, result := LoginRequired(r)
	encoder := json.NewEncoder(w)
	if !result {
		http.Error(w, "not logined", 401)
		return
	}
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	project := projectModel.FindProject(projectID)
	if project == nil && project.UserId.Int64 != current_user.Id {
		http.Error(w, "project not found", 404)
		return
	}
	encoder.Encode(project)
	return
}

func (u *Projects) Create(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, result := LoginRequired(r)
	encoder := json.NewEncoder(w)
	if !result {
		http.Error(w, "not logined", 401)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Wrong Form", 400)
		return
	}
	var newProjectForm NewProjectForm
	err = param.Parse(r.PostForm, &newProjectForm)

	if err != nil {
		http.Error(w, "Wrong parameter", 500)
		return
	}
	fmt.Printf("post new project parameter: %+v\n", newProjectForm)
	project := projectModel.NewProject(0, current_user.Id, newProjectForm.Title, newProjectForm.Description)
	if !project.Save() {
		http.Error(w, "save failed", 500)
		return
	}

	// 初期リストを作っておく
	todo := listModel.NewList(0, project.Id, current_user.Id, "ToDo", "ff0000")
	inprogress := listModel.NewList(0, project.Id, current_user.Id, "InProgress", "0000ff")
	done := listModel.NewList(0, project.Id, current_user.Id, "Done", "0a0a0a")
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

	encoder.Encode(*project)
}

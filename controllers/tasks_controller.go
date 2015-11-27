package controllers

import (
	listModel "../models/list"
	projectModel "../models/project"
	taskModel "../models/task"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/goji/param"
	"github.com/zenazn/goji/web"
	"net/http"
	"strconv"
)

type Tasks struct {
}

type NewTaskForm struct {
	Title string `param:"title"`
}

type MoveTaskFrom struct {
	ToListId     int64 `param:"to_list_id"`
	PrevToTaskId int64 `param:"prev_to_task_id"`
}

func (u *Tasks) Index(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, result := LoginRequired(r)
	encoder := json.NewEncoder(w)
	if !result {
		http.Error(w, "not logined", 401)
		return
	}
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	parentProject := projectModel.FindProject(projectID)
	if parentProject == nil && parentProject.UserId.Int64 != current_user.Id {
		http.Error(w, "project not found", 404)
		return
	}
	listID, _ := strconv.ParseInt(c.URLParams["list_id"], 10, 64)
	parentList := listModel.FindList(projectID, listID)
	if parentList == nil {
		http.Error(w, "list not found", 404)
		return
	}
	tasks := parentList.Tasks()
	encoder.Encode(tasks)
	return
}

func (u *Tasks) Create(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, result := LoginRequired(r)
	encoder := json.NewEncoder(w)
	if !result {
		http.Error(w, "not logined", 401)
		return
	}
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	parentProject := projectModel.FindProject(projectID)
	if parentProject == nil && parentProject.UserId.Int64 != current_user.Id {
		http.Error(w, "project not found", 404)
		return
	}
	listID, _ := strconv.ParseInt(c.URLParams["list_id"], 10, 64)
	parentList := listModel.FindList(projectID, listID)
	if parentList == nil {
		http.Error(w, "list not found", 404)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Wrong From", 400)
		return
	}
	var newTaskForm NewTaskForm
	err = param.Parse(r.PostForm, &newTaskForm)
	if err != nil {
		http.Error(w, "Wrong parameter", 500)
		return
	}
	fmt.Printf("post new task parameter: %+v\n", newTaskForm)

	task := taskModel.NewTask(0, parentList.Id, parentList.UserId, sql.NullInt64{}, newTaskForm.Title)

	repo := parentProject.Repository()
	if !task.Save(repo, &current_user.OauthToken) {
		http.Error(w, "save failed", 500)
		return
	}
	encoder.Encode(*task)
}

func (u *Tasks) MoveTask(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, result := LoginRequired(r)
	encoder := json.NewEncoder(w)
	if !result {
		http.Error(w, "not logined", 401)
		return
	}
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	parentProject := projectModel.FindProject(projectID)
	if parentProject == nil || parentProject.UserId.Int64 != current_user.Id {
		http.Error(w, "project not found", 404)
		return
	}

	listID, _ := strconv.ParseInt(c.URLParams["list_id"], 10, 64)
	parentList := listModel.FindList(parentProject.Id, listID)
	if parentList == nil {
		http.Error(w, "list not found", 404)
		return
	}

	taskID, _ := strconv.ParseInt(c.URLParams["task_id"], 10, 64)
	task := taskModel.FindTask(parentList.Id, taskID)

	err := r.ParseForm()
	if err != nil {
		fmt.Printf("err: %+v\n", err)
		http.Error(w, "WrongForm", 400)
		return
	}
	var moveTaskFrom MoveTaskFrom
	err = param.Parse(r.PostForm, &moveTaskFrom)
	if err != nil {
		fmt.Printf("err: %+v\n", err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	fmt.Printf("post move taks parameter: %+v\n", moveTaskFrom)
	var prevToTaskId *int64
	if moveTaskFrom.PrevToTaskId != 0 {
		prevToTaskId = &moveTaskFrom.PrevToTaskId
	}

	repo := parentProject.Repository()
	if !task.ChangeList(moveTaskFrom.ToListId, prevToTaskId, repo, &current_user.OauthToken) {
		http.Error(w, "failed change list", 500)
		return
	}
	allLists := parentProject.Lists()
	for _, l := range allLists {
		l.ListTasks = l.Tasks()
	}
	fmt.Printf("move task: %+v\n", allLists)
	encoder.Encode(allLists)
	return
}

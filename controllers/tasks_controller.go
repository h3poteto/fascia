package controllers

import (
	listModel "../models/list"
	projectModel "../models/project"
	taskModel "../models/task"
	"../modules/logging"
	"database/sql"
	"encoding/json"
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
	current_user, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Index").Errorf("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	parentProject := projectModel.FindProject(projectID)
	if parentProject == nil || parentProject.UserId.Int64 != current_user.Id {
		logging.SharedInstance().MethodInfo("TasksController", "Index").Error("project not found")
		http.Error(w, "project not found", 404)
		return
	}
	listID, _ := strconv.ParseInt(c.URLParams["list_id"], 10, 64)
	parentList := listModel.FindList(projectID, listID)
	if parentList == nil {
		logging.SharedInstance().MethodInfo("TasksController", "Index").Error("list not found")
		http.Error(w, "list not found", 404)
		return
	}
	tasks := parentList.Tasks()
	encoder.Encode(tasks)
	return
}

func (u *Tasks) Create(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Create").Errorf("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	parentProject := projectModel.FindProject(projectID)
	if parentProject == nil || parentProject.UserId.Int64 != current_user.Id {
		logging.SharedInstance().MethodInfo("TasksController", "Create").Error("project not found")
		http.Error(w, "project not found", 404)
		return
	}
	listID, _ := strconv.ParseInt(c.URLParams["list_id"], 10, 64)
	parentList := listModel.FindList(projectID, listID)
	if parentList == nil {
		logging.SharedInstance().MethodInfo("TasksController", "Create").Error("list not found")
		http.Error(w, "list not found", 404)
		return
	}

	err = r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Create").Errorf("wrong form: %v", err)
		http.Error(w, "Wrong Form", 400)
		return
	}
	var newTaskForm NewTaskForm
	err = param.Parse(r.PostForm, &newTaskForm)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Create").Errorf("wrong parameter: %v", err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("TasksController", "Create").Debugf("post new task parameter: %+v", newTaskForm)

	task := taskModel.NewTask(0, parentList.Id, parentList.UserId, sql.NullInt64{}, newTaskForm.Title)

	repo := parentProject.Repository()
	if !task.Save(repo, &current_user.OauthToken) {
		logging.SharedInstance().MethodInfo("TasksController", "Create").Error("save failed")
		http.Error(w, "save failed", 500)
		return
	}
	logging.SharedInstance().MethodInfo("TasksController", "Create").Info("success to create task")
	encoder.Encode(*task)
}

func (u *Tasks) MoveTask(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask").Errorf("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, _ := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	parentProject := projectModel.FindProject(projectID)
	if parentProject == nil || parentProject.UserId.Int64 != current_user.Id {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask").Error("project not found")
		http.Error(w, "project not found", 404)
		return
	}

	listID, _ := strconv.ParseInt(c.URLParams["list_id"], 10, 64)
	parentList := listModel.FindList(parentProject.Id, listID)
	if parentList == nil {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask").Error("list not found")
		http.Error(w, "list not found", 404)
		return
	}

	taskID, _ := strconv.ParseInt(c.URLParams["task_id"], 10, 64)
	task := taskModel.FindTask(parentList.Id, taskID)

	err = r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask").Errorf("wrong form: %v", err)
		http.Error(w, "Wrong Form", 400)
		return
	}
	var moveTaskFrom MoveTaskFrom
	err = param.Parse(r.PostForm, &moveTaskFrom)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask").Errorf("wrong parameter: %v", err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("TasksController", "MoveTask").Debugf("post move taks parameter: %+v", moveTaskFrom)
	var prevToTaskId *int64
	if moveTaskFrom.PrevToTaskId != 0 {
		prevToTaskId = &moveTaskFrom.PrevToTaskId
	}

	repo := parentProject.Repository()
	if !task.ChangeList(moveTaskFrom.ToListId, prevToTaskId, repo, &current_user.OauthToken) {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask").Error("failed change list")
		http.Error(w, "failed change list", 500)
		return
	}
	allLists := parentProject.Lists()
	for _, l := range allLists {
		l.ListTasks = l.Tasks()
	}
	logging.SharedInstance().MethodInfo("TasksController", "MoveTask").Debugf("move task: %+v", allLists)
	logging.SharedInstance().MethodInfo("TasksController", "MoveTask").Info("success to move task")
	encoder.Encode(allLists)
	return
}

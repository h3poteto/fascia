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
	Title       string `param:"title"`
	Description string `param:"description"`
}

type MoveTaskFrom struct {
	ToListId     int64 `param:"to_list_id"`
	PrevToTaskId int64 `param:"prev_to_task_id"`
}

type TaskJsonFormat struct {
	Id          int64
	ListId      int64
	UserId      int64
	IssueNumber int64
	Title       string
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
	if parentProject == nil || parentProject.UserId != current_user.Id {
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
	jsonTasks := make([]*TaskJsonFormat, 0)
	for _, t := range tasks {
		jsonTasks = append(jsonTasks, &TaskJsonFormat{Id: t.Id, ListId: t.ListId, UserId: t.UserId, IssueNumber: t.IssueNumber.Int64, Title: t.Title})
	}
	encoder.Encode(jsonTasks)
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
	if parentProject == nil || parentProject.UserId != current_user.Id {
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

	task := taskModel.NewTask(0, parentList.Id, parentList.UserId, sql.NullInt64{}, newTaskForm.Title, newTaskForm.Description)

	repo := parentProject.Repository()
	if !task.Save(repo, &current_user.OauthToken) {
		logging.SharedInstance().MethodInfo("TasksController", "Create").Error("save failed")
		http.Error(w, "save failed", 500)
		return
	}
	jsonTask := TaskJsonFormat{Id: task.Id, ListId: task.ListId, UserId: task.UserId, IssueNumber: task.IssueNumber.Int64, Title: task.Title}
	logging.SharedInstance().MethodInfo("TasksController", "Create").Info("success to create task")
	encoder.Encode(jsonTask)
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
	if parentProject == nil || parentProject.UserId != current_user.Id {
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
	jsonLists := make([]*ListJSONFormat, 0)
	for _, l := range allLists {
		jsonLists = append(jsonLists, &ListJSONFormat{Id: l.Id, ProjectId: l.ProjectId, UserId: l.UserId, Title: l.Title.String, ListTasks: TaskFormatToJson(l.Tasks()), Color: l.Color.String, ListOptionId: l.ListOptionId.Int64})
	}
	noneList := parentProject.NoneList()
	jsonNoneList := &ListJSONFormat{Id: noneList.Id, ProjectId: noneList.ProjectId, UserId: noneList.UserId, Title: noneList.Title.String, ListTasks: TaskFormatToJson(noneList.Tasks()), Color: noneList.Color.String, ListOptionId: noneList.ListOptionId.Int64}
	jsonAllLists := AllListJSONFormat{Lists: jsonLists, NoneList: jsonNoneList}
	logging.SharedInstance().MethodInfo("TasksController", "MoveTask").Debugf("move task: %+v", allLists)
	logging.SharedInstance().MethodInfo("TasksController", "MoveTask").Info("success to move task")
	encoder.Encode(jsonAllLists)
	return
}

func TaskFormatToJson(tasks []*taskModel.TaskStruct) []*TaskJsonFormat {
	jsonTasks := make([]*TaskJsonFormat, 0)
	for _, t := range tasks {
		jsonTasks = append(jsonTasks, &TaskJsonFormat{Id: t.Id, ListId: t.ListId, UserId: t.UserId, IssueNumber: t.IssueNumber.Int64, Title: t.Title})
	}
	return jsonTasks
}

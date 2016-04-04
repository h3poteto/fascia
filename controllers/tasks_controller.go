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
	ToListID     int64 `param:"to_list_id"`
	PrevToTaskID int64 `param:"prev_to_task_id"`
}

type TaskJsonFormat struct {
	ID          int64
	ListID      int64
	UserID      int64
	IssueNumber int64
	Title       string
	PullRequest bool
}

func (u *Tasks) Index(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Index", false, c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Index", true, c).Errorf("parse error: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	parentProject, err := projectModel.FindProject(projectID)
	if err != nil || parentProject.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("TasksController", "Index", false, c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	listID, err := strconv.ParseInt(c.URLParams["list_id"], 10, 64)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Index", true, c).Errorf("parse error: %v", err)
		http.Error(w, "list not found", 404)
		return
	}
	parentList, err := listModel.FindList(projectID, listID)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Index", false, c).Warnf("list not found: %v", err)
		http.Error(w, "list not found", 404)
		return
	}
	tasks := parentList.Tasks()
	jsonTasks := make([]*TaskJsonFormat, 0)
	for _, t := range tasks {
		jsonTasks = append(jsonTasks, &TaskJsonFormat{ID: t.ID, ListID: t.ListID, UserID: t.UserID, IssueNumber: t.IssueNumber.Int64, Title: t.Title, PullRequest: t.PullRequest})
	}
	encoder.Encode(jsonTasks)
	logging.SharedInstance().MethodInfo("TasksController", "Index", false, c).Info("success to get tasks")
	return
}

func (u *Tasks) Create(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Create", false, c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Create", true, c).Errorf("parse error: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	parentProject, err := projectModel.FindProject(projectID)
	if err != nil || parentProject.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("TasksController", "Create", false, c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	listID, err := strconv.ParseInt(c.URLParams["list_id"], 10, 64)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Create", true, c).Errorf("parse error: %v", err)
		http.Error(w, "list not found", 404)
		return
	}
	parentList, err := listModel.FindList(projectID, listID)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Create", false, c).Warnf("list not found: %v", err)
		http.Error(w, "list not found", 404)
		return
	}

	err = r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Create", true, c).Errorf("wrong form: %v", err)
		http.Error(w, "Wrong Form", 400)
		return
	}
	var newTaskForm NewTaskForm
	err = param.Parse(r.PostForm, &newTaskForm)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Create", true, c).Errorf("wrong parameter: %v", err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("TasksController", "Create", false, c).Debugf("post new task parameter: %+v", newTaskForm)

	task := taskModel.NewTask(0, parentList.ID, parentProject.ID, parentList.UserID, sql.NullInt64{}, newTaskForm.Title, newTaskForm.Description, false, sql.NullString{})

	repo, _ := parentProject.Repository()
	if err := task.Save(repo, &currentUser.OauthToken); err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Create", true, c).Errorf("save failed: %v", err)
		http.Error(w, "save failed", 500)
		return
	}
	jsonTask := TaskJsonFormat{ID: task.ID, ListID: task.ListID, UserID: task.UserID, IssueNumber: task.IssueNumber.Int64, Title: task.Title, PullRequest: task.PullRequest}
	logging.SharedInstance().MethodInfo("TasksController", "Create", false, c).Info("success to create task")
	encoder.Encode(jsonTask)
}

func (u *Tasks) MoveTask(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask", false, c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask", true, c).Errorf("parse error: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	parentProject, err := projectModel.FindProject(projectID)
	if err != nil || parentProject.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask", false, c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}

	listID, err := strconv.ParseInt(c.URLParams["list_id"], 10, 64)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask", true, c).Errorf("parse error: %v", err)
		http.Error(w, "list not found", 404)
		return
	}
	parentList, err := listModel.FindList(parentProject.ID, listID)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask", false, c).Warnf("list not found: %v", err)
		http.Error(w, "list not found", 404)
		return
	}

	taskID, err := strconv.ParseInt(c.URLParams["task_id"], 10, 64)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask", true, c).Errorf("parse error: %v", err)
		http.Error(w, "task not found", 404)
		return
	}
	task, err := taskModel.FindTask(parentList.ID, taskID)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask", true, c).Errorf("find task error: %v", err)
		http.Error(w, "task not find", 500)
		return
	}

	err = r.ParseForm()
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask", true, c).Errorf("wrong form: %v", err)
		http.Error(w, "Wrong Form", 400)
		return
	}
	var moveTaskFrom MoveTaskFrom
	err = param.Parse(r.PostForm, &moveTaskFrom)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask", true, c).Errorf("wrong parameter: %v", err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("TasksController", "MoveTask", false, c).Debugf("post move taks parameter: %+v", moveTaskFrom)
	var prevToTaskID *int64
	if moveTaskFrom.PrevToTaskID != 0 {
		prevToTaskID = &moveTaskFrom.PrevToTaskID
	}

	repo, _ := parentProject.Repository()
	if err := task.ChangeList(moveTaskFrom.ToListID, prevToTaskID, repo, &currentUser.OauthToken); err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask", true, c).Errorf("failed change list: %v", err)
		http.Error(w, "failed change list", 500)
		return
	}
	allLists := parentProject.Lists()
	var jsonLists []*ListJSONFormat
	for _, l := range allLists {
		jsonLists = append(jsonLists, &ListJSONFormat{ID: l.ID, ProjectID: l.ProjectID, UserID: l.UserID, Title: l.Title.String, ListTasks: TaskFormatToJson(l.Tasks()), Color: l.Color.String, ListOptionID: l.ListOptionID.Int64})
	}
	noneList, err := parentProject.NoneList()
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask", true, c).Error(err)
		http.Error(w, "none list not found", 500)
		return
	}
	jsonNoneList := &ListJSONFormat{ID: noneList.ID, ProjectID: noneList.ProjectID, UserID: noneList.UserID, Title: noneList.Title.String, ListTasks: TaskFormatToJson(noneList.Tasks()), Color: noneList.Color.String, ListOptionID: noneList.ListOptionID.Int64}
	jsonAllLists := AllListJSONFormat{Lists: jsonLists, NoneList: jsonNoneList}
	logging.SharedInstance().MethodInfo("TasksController", "MoveTask", false, c).Debugf("move task: %+v", allLists)
	logging.SharedInstance().MethodInfo("TasksController", "MoveTask", false, c).Info("success to move task")
	encoder.Encode(jsonAllLists)
	return
}

func TaskFormatToJson(tasks []*taskModel.TaskStruct) []*TaskJsonFormat {
	jsonTasks := make([]*TaskJsonFormat, 0)
	for _, t := range tasks {
		jsonTasks = append(jsonTasks, &TaskJsonFormat{ID: t.ID, ListID: t.ListID, UserID: t.UserID, IssueNumber: t.IssueNumber.Int64, Title: t.Title, PullRequest: t.PullRequest})
	}
	return jsonTasks
}

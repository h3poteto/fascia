package controllers

import (
	listModel "../models/list"
	projectModel "../models/project"
	taskModel "../models/task"
	"../modules/logging"

	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/goji/param"
	"github.com/pkg/errors"
	"github.com/zenazn/goji/web"
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

type TaskJSONFormat struct {
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
		logging.SharedInstance().MethodInfo("TasksController", "Index", c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Index", err, c).Error(err)
		http.Error(w, "project not found", 404)
		return
	}
	parentProject, err := projectModel.FindProject(projectID)
	if err != nil || parentProject.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("TasksController", "Index", c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	listID, err := strconv.ParseInt(c.URLParams["list_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Index", err, c).Error(err)
		http.Error(w, "list not found", 404)
		return
	}
	parentList, err := listModel.FindList(projectID, listID)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Index", c).Warnf("list not found: %v", err)
		http.Error(w, "list not found", 404)
		return
	}
	tasks, err := parentList.Tasks()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Index", err, c).Error(err)
		http.Error(w, "task not found", 500)
		return
	}
	jsonTasks := make([]*TaskJSONFormat, 0)
	for _, t := range tasks {
		jsonTasks = append(jsonTasks, &TaskJSONFormat{ID: t.ID, ListID: t.ListID, UserID: t.UserID, IssueNumber: t.IssueNumber.Int64, Title: t.Title, PullRequest: t.PullRequest})
	}
	encoder.Encode(jsonTasks)
	logging.SharedInstance().MethodInfo("TasksController", "Index", c).Info("success to get tasks")
	return
}

func (u *Tasks) Create(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Create", c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Create", err, c).Error(err)
		http.Error(w, "project not found", 404)
		return
	}
	parentProject, err := projectModel.FindProject(projectID)
	if err != nil || parentProject.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("TasksController", "Create", c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}
	listID, err := strconv.ParseInt(c.URLParams["list_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Create", err, c).Error(err)
		http.Error(w, "list not found", 404)
		return
	}
	parentList, err := listModel.FindList(projectID, listID)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Create", c).Warnf("list not found: %v", err)
		http.Error(w, "list not found", 404)
		return
	}

	err = r.ParseForm()
	if err != nil {
		err := errors.Wrap(err, "wrong form")
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Create", err, c).Error(err)
		http.Error(w, "Wrong Form", 400)
		return
	}
	var newTaskForm NewTaskForm
	err = param.Parse(r.PostForm, &newTaskForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Create", err, c).Error(err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("TasksController", "Create", c).Debugf("post new task parameter: %+v", newTaskForm)

	task := taskModel.NewTask(0, parentList.ID, parentProject.ID, parentList.UserID, sql.NullInt64{}, newTaskForm.Title, newTaskForm.Description, false, sql.NullString{})

	repo, _ := parentProject.Repository()
	if err := task.Save(repo, &currentUser.OauthToken); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Create", err, c).Errorf("save failed: %v", err)
		http.Error(w, "save failed", 500)
		return
	}
	jsonTask := TaskJSONFormat{ID: task.ID, ListID: task.ListID, UserID: task.UserID, IssueNumber: task.IssueNumber.Int64, Title: task.Title, PullRequest: task.PullRequest}
	logging.SharedInstance().MethodInfo("TasksController", "Create", c).Info("success to create task")
	encoder.Encode(jsonTask)
}

func (u *Tasks) MoveTask(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask", c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}
	encoder := json.NewEncoder(w)
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "MoveTask", err, c).Error(err)
		http.Error(w, "project not found", 404)
		return
	}
	parentProject, err := projectModel.FindProject(projectID)
	if err != nil || parentProject.UserID != currentUser.ID {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask", c).Warnf("project not found: %v", err)
		http.Error(w, "project not found", 404)
		return
	}

	listID, err := strconv.ParseInt(c.URLParams["list_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "MoveTask", err, c).Error(err)
		http.Error(w, "list not found", 404)
		return
	}
	parentList, err := listModel.FindList(parentProject.ID, listID)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask", c).Warnf("list not found: %v", err)
		http.Error(w, "list not found", 404)
		return
	}

	taskID, err := strconv.ParseInt(c.URLParams["task_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "MoveTask", err, c).Error(err)
		http.Error(w, "task not found", 404)
		return
	}
	task, err := taskModel.FindTask(parentList.ID, taskID)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "MoveTask", err, c).Errorf("find task error: %v", err)
		http.Error(w, "task not find", 500)
		return
	}

	err = r.ParseForm()
	if err != nil {
		err := errors.Wrap(err, "wrong form")
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "MoveTask", err, c).Error(err)
		http.Error(w, "Wrong Form", 400)
		return
	}
	var moveTaskFrom MoveTaskFrom
	err = param.Parse(r.PostForm, &moveTaskFrom)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "MoveTask", err, c).Error(err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("TasksController", "MoveTask", c).Debugf("post move taks parameter: %+v", moveTaskFrom)
	var prevToTaskID *int64
	if moveTaskFrom.PrevToTaskID != 0 {
		prevToTaskID = &moveTaskFrom.PrevToTaskID
	}

	repo, _ := parentProject.Repository()
	if err := task.ChangeList(moveTaskFrom.ToListID, prevToTaskID, repo, &currentUser.OauthToken); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "MoveTask", err, c).Errorf("failed change list: %v", err)
		http.Error(w, "failed change list", 500)
		return
	}
	allLists, err := parentProject.Lists()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "MoveTask", err, c).Error(err)
		http.Error(w, "lists not found", 500)
		return
	}
	jsonLists := ListsFormatToJSON(allLists)
	noneList, err := parentProject.NoneList()
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "MoveTask", err, c).Error(err)
		http.Error(w, "none list not found", 500)
		return
	}
	jsonNoneList := ListFormatToJSON(noneList)
	jsonAllLists := AllListJSONFormat{Lists: jsonLists, NoneList: jsonNoneList}
	logging.SharedInstance().MethodInfo("TasksController", "MoveTask", c).Debugf("move task: %+v", allLists)
	logging.SharedInstance().MethodInfo("TasksController", "MoveTask", c).Info("success to move task")
	encoder.Encode(jsonAllLists)
	return
}

// TaskFormatToJSON convert task model's array to json
func TaskFormatToJSON(tasks []*taskModel.TaskStruct) []*TaskJSONFormat {
	jsonTasks := make([]*TaskJSONFormat, 0)
	for _, t := range tasks {
		jsonTasks = append(jsonTasks, &TaskJSONFormat{ID: t.ID, ListID: t.ListID, UserID: t.UserID, IssueNumber: t.IssueNumber.Int64, Title: t.Title, PullRequest: t.PullRequest})
	}
	return jsonTasks
}

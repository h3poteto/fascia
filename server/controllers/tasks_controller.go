package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/services"
	"github.com/h3poteto/fascia/server/validators"
	"github.com/h3poteto/fascia/server/views"

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

// MoveTaskForm is struct for move task
type MoveTaskForm struct {
	ToListID     int64 `param:"to_list_id"`
	PrevToTaskID int64 `param:"prev_to_task_id"`
}

// EditTaskForm is struct for edit task
type EditTaskForm struct {
	Title       string `param:"title"`
	Description string `param:"description"`
}

func (u *Tasks) Create(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Create", c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}

	projectService, parentList, statusCode, err := setProjectAndList(c, w, currentUser)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Create", err, c).Error(err)
		switch statusCode {
		case 404:
			http.Error(w, "Not Found", 404)
		default:
			http.Error(w, "Internal Server Error", 500)
		}
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

	valid, err := validators.TaskCreateValidation(newTaskForm.Title, newTaskForm.Description)
	if err != nil || !valid {
		logging.SharedInstance().MethodInfo("TasksController", "Create", c).Infof("validation error: %v", err)
		http.Error(w, "validation error", 422)
		return
	}

	task := services.NewTask(
		0,
		parentList.ListEntity.ListModel.ID,
		projectService.ProjectEntity.ProjectModel.ID,
		parentList.ListEntity.ListModel.UserID,
		sql.NullInt64{},
		newTaskForm.Title,
		newTaskForm.Description,
		false,
		sql.NullString{},
	)

	if err := task.Save(); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Create", err, c).Errorf("save failed: %v", err)
		http.Error(w, "save failed", 500)
		return
	}

	encoder := json.NewEncoder(w)

	jsonAllLists, err := allListsResponse(projectService)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	encoder.Encode(jsonAllLists)
	logging.SharedInstance().MethodInfo("TasksController", "Create", c).Debugf("create task success: %+v", task)
	logging.SharedInstance().MethodInfo("TasksController", "Create", c).Info("success to create task")
	return
}

// Show render json with task detail
func (u *Tasks) Show(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Show", c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}

	_, parentList, statusCode, err := setProjectAndList(c, w, currentUser)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Show", err, c).Error(err)
		switch statusCode {
		case 404:
			http.Error(w, "Not Found", 404)
		default:
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}

	task, statusCode, err := setTask(c, w, parentList)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Show", err, c).Error(err)
		switch statusCode {
		case 404:
			http.Error(w, "Not Found", 404)
		default:
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}

	encoder := json.NewEncoder(w)
	jsonTask, err := views.ParseTaskJSON(task.TaskEntity)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Show", err, c).Error(err)
		http.Error(w, "task error", 500)
		return
	}
	logging.SharedInstance().MethodInfo("TasksController", "Show", c).Info("success to get task")
	encoder.Encode(jsonTask)
	return
}

func (u *Tasks) MoveTask(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask", c).Infof("login error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}

	projectService, parentList, statusCode, err := setProjectAndList(c, w, currentUser)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "MoveTask", err, c).Error(err)
		switch statusCode {
		case 404:
			http.Error(w, "Not Found", 404)
		default:
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}

	task, statusCode, err := setTask(c, w, parentList)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "MoveTask", err, c).Error(err)
		switch statusCode {
		case 404:
			http.Error(w, "Not Found", 404)
		default:
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}

	err = r.ParseForm()
	if err != nil {
		err := errors.Wrap(err, "wrong form")
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "MoveTask", err, c).Error(err)
		http.Error(w, "Wrong Form", 400)
		return
	}
	var moveTaskFrom MoveTaskForm
	err = param.Parse(r.PostForm, &moveTaskFrom)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "MoveTask", err, c).Error(err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("TasksController", "MoveTask", c).Debugf("post move taks parameter: %+v", moveTaskFrom)

	valid, err := validators.TaskMoveValidation(moveTaskFrom.ToListID, moveTaskFrom.PrevToTaskID)
	if err != nil || !valid {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask", c).Infof("validation error: %v", err)
		http.Error(w, "validation error", 422)
		return
	}

	var prevToTaskID *int64
	if moveTaskFrom.PrevToTaskID != 0 {
		prevToTaskID = &moveTaskFrom.PrevToTaskID
	}

	if err := task.ChangeList(moveTaskFrom.ToListID, prevToTaskID); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "MoveTask", err, c).Errorf("failed change list: %v", err)
		http.Error(w, "failed change list", 500)
		return
	}

	encoder := json.NewEncoder(w)

	jsonAllLists, err := allListsResponse(projectService)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	encoder.Encode(jsonAllLists)

	logging.SharedInstance().MethodInfo("TasksController", "MoveTask", c).Debugf("move task: %+v", task)
	logging.SharedInstance().MethodInfo("TasksController", "MoveTask", c).Info("success to move task")
	return
}

// Update a task
func (u *Tasks) Update(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Update", c).Infof("loging error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}

	projectService, parentList, statusCode, err := setProjectAndList(c, w, currentUser)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Update", err, c).Error(err)
		switch statusCode {
		case 404:
			http.Error(w, "Not Found", 404)
		default:
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}

	task, statusCode, err := setTask(c, w, parentList)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Update", err, c).Error(err)
		switch statusCode {
		case 404:
			http.Error(w, "Not Found", 404)
		default:
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}

	err = r.ParseForm()
	if err != nil {
		err := errors.Wrap(err, "wrong form")
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Update", err, c).Error(err)
		http.Error(w, "Wrong Form", 400)
		return
	}
	var editTaskForm EditTaskForm
	err = param.Parse(r.PostForm, &editTaskForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Update", err, c).Error(err)
		http.Error(w, "Wrong parameter", 500)
		return
	}
	logging.SharedInstance().MethodInfo("TasksController", "Update", c).Debugf("post update parameter: %+v", editTaskForm)

	valid, err := validators.TaskUpdateValidation(editTaskForm.Title, editTaskForm.Description)
	if err != nil || !valid {
		logging.SharedInstance().MethodInfo("TasksController", "Update", c).Infof("validation error: %v", err)
		http.Error(w, "validation error", 422)
		return
	}

	err = task.Update(
		task.TaskEntity.TaskModel.ListID,
		task.TaskEntity.TaskModel.IssueNumber,
		editTaskForm.Title,
		editTaskForm.Description,
		task.TaskEntity.TaskModel.PullRequest,
		task.TaskEntity.TaskModel.HTMLURL,
	)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Update", err, c).Error(err)
		http.Error(w, "update error", 500)
		return
	}

	encoder := json.NewEncoder(w)

	jsonAllLists, err := allListsResponse(projectService)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	encoder.Encode(jsonAllLists)
	logging.SharedInstance().MethodInfo("TasksController", "Update", c).Debugf("update task: %+v", task)
	logging.SharedInstance().MethodInfo("TasksController", "Update", c).Info("success to update task")
	return
}

// Delete a task
func (u *Tasks) Delete(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUser, err := LoginRequired(r)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Delete", c).Infof("loging error: %v", err)
		http.Error(w, "not logined", 401)
		return
	}

	projectService, parentList, statusCode, err := setProjectAndList(c, w, currentUser)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Delete", err, c).Error(err)
		switch statusCode {
		case 404:
			http.Error(w, "Not Found", 404)
		default:
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}

	task, statusCode, err := setTask(c, w, parentList)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Delete", err, c).Error(err)
		switch statusCode {
		case 404:
			http.Error(w, "Not Found", 404)
		default:
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}

	err = task.Delete()
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Delete", c).Info(err)
		http.Error(w, "Bad Request", 400)
		return
	}
	encoder := json.NewEncoder(w)
	jsonAllLists, err := allListsResponse(projectService)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	encoder.Encode(jsonAllLists)
	logging.SharedInstance().MethodInfo("TasksController", "Delete", c).Info("success to delete a task")
	return
}

func setProjectAndList(c web.C, w http.ResponseWriter, currentUser *services.User) (*services.Project, *services.List, int, error) {
	projectID, err := strconv.ParseInt(c.URLParams["project_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		return nil, nil, 404, err
	}
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
		return nil, nil, 404, err
	}
	listID, err := strconv.ParseInt(c.URLParams["list_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		return nil, nil, 404, err
	}
	parentList, err := handlers.FindList(projectService.ProjectEntity.ProjectModel.ID, listID)
	if err != nil {
		return nil, nil, 404, err
	}
	return projectService, parentList, 200, nil
}

func setTask(c web.C, w http.ResponseWriter, list *services.List) (*services.Task, int, error) {
	taskID, err := strconv.ParseInt(c.URLParams["task_id"], 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		return nil, 404, err
	}
	task, err := handlers.FindTask(list.ListEntity.ListModel.ID, taskID)
	if err != nil {
		return nil, 404, err
	}

	return task, 200, nil
}

func allListsResponse(projectService *services.Project) (*views.AllLists, error) {
	allLists, err := projectService.ProjectEntity.Lists()
	if err != nil {
		return nil, err
	}
	noneList, err := projectService.ProjectEntity.NoneList()
	if err != nil {
		return nil, err
	}
	return views.ParseAllListsJSON(noneList, allLists)
}

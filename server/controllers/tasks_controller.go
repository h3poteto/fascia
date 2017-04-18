package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/services"
	"github.com/h3poteto/fascia/server/validators"
	"github.com/h3poteto/fascia/server/views"

	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
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

func (u *Tasks) Create(c echo.Context) error {
	currentUser, err := LoginRequired(c)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Create", c).Infof("login error: %v", err)
		return NewJSONError(err, http.StatusUnauthorized, c)
	}

	projectService, parentList, statusCode, err := setProjectAndList(c, currentUser)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Create", err, c).Error(err)
		switch statusCode {
		case 404:
			return NewJSONError(err, http.StatusNotFound, c)
		default:
			return err
		}
	}

	var newTaskForm NewTaskForm
	err = c.Bind(newTaskForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Create", err, c).Error(err)
		return err
	}
	logging.SharedInstance().MethodInfo("TasksController", "Create", c).Debugf("post new task parameter: %+v", newTaskForm)

	valid, err := validators.TaskCreateValidation(newTaskForm.Title, newTaskForm.Description)
	if err != nil || !valid {
		logging.SharedInstance().MethodInfo("TasksController", "Create", c).Infof("validation error: %v", err)
		return NewJSONError(err, http.StatusUnprocessableEntity, c)
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
		return err
	}

	jsonAllLists, err := allListsResponse(projectService)
	if err != nil {
		return err
	}
	logging.SharedInstance().MethodInfo("TasksController", "Create", c).Debugf("create task success: %+v", task)
	logging.SharedInstance().MethodInfo("TasksController", "Create", c).Info("success to create task")
	return c.JSON(http.StatusOK, jsonAllLists)
}

// Show render json with task detail
func (u *Tasks) Show(c echo.Context) error {
	currentUser, err := LoginRequired(c)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Show", c).Infof("login error: %v", err)
		return NewJSONError(err, http.StatusUnauthorized, c)
	}

	_, parentList, statusCode, err := setProjectAndList(c, currentUser)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Show", err, c).Error(err)
		switch statusCode {
		case 404:
			return NewJSONError(err, http.StatusNotFound, c)
		default:
			return err
		}
	}

	task, statusCode, err := setTask(c, parentList)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Show", err, c).Error(err)
		switch statusCode {
		case 404:
			return NewJSONError(err, http.StatusNotFound, c)
		default:
			return err
		}
	}

	jsonTask, err := views.ParseTaskJSON(task.TaskEntity)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Show", err, c).Error(err)
		return err
	}
	logging.SharedInstance().MethodInfo("TasksController", "Show", c).Info("success to get task")
	return c.JSON(http.StatusOK, jsonTask)
}

func (u *Tasks) MoveTask(c echo.Context) error {
	currentUser, err := LoginRequired(c)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask", c).Infof("login error: %v", err)
		return NewJSONError(err, http.StatusUnauthorized, c)
	}

	projectService, parentList, statusCode, err := setProjectAndList(c, currentUser)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "MoveTask", err, c).Error(err)
		switch statusCode {
		case 404:
			return NewJSONError(err, http.StatusNotFound, c)
		default:
			return err
		}
	}

	task, statusCode, err := setTask(c, parentList)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "MoveTask", err, c).Error(err)
		switch statusCode {
		case 404:
			return NewJSONError(err, http.StatusNotFound, c)
		default:
			return err
		}
	}

	var moveTaskFrom MoveTaskForm
	err = c.Bind(moveTaskFrom)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "MoveTask", err, c).Error(err)
		return err
	}
	logging.SharedInstance().MethodInfo("TasksController", "MoveTask", c).Debugf("post move taks parameter: %+v", moveTaskFrom)

	valid, err := validators.TaskMoveValidation(moveTaskFrom.ToListID, moveTaskFrom.PrevToTaskID)
	if err != nil || !valid {
		logging.SharedInstance().MethodInfo("TasksController", "MoveTask", c).Infof("validation error: %v", err)
		return NewJSONError(err, http.StatusUnprocessableEntity, c)
	}

	var prevToTaskID *int64
	if moveTaskFrom.PrevToTaskID != 0 {
		prevToTaskID = &moveTaskFrom.PrevToTaskID
	}

	if err := task.ChangeList(moveTaskFrom.ToListID, prevToTaskID); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "MoveTask", err, c).Errorf("failed change list: %v", err)
		return err
	}

	jsonAllLists, err := allListsResponse(projectService)
	if err != nil {
		return err
	}
	logging.SharedInstance().MethodInfo("TasksController", "MoveTask", c).Debugf("move task: %+v", task)
	logging.SharedInstance().MethodInfo("TasksController", "MoveTask", c).Info("success to move task")
	return c.JSON(http.StatusOK, jsonAllLists)
}

// Update a task
func (u *Tasks) Update(c echo.Context) error {
	currentUser, err := LoginRequired(c)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Update", c).Infof("loging error: %v", err)
		return NewJSONError(err, http.StatusUnauthorized, c)
	}

	projectService, parentList, statusCode, err := setProjectAndList(c, currentUser)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Update", err, c).Error(err)
		switch statusCode {
		case 404:
			return NewJSONError(err, http.StatusNotFound, c)
		default:
			return err
		}
	}

	task, statusCode, err := setTask(c, parentList)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Update", err, c).Error(err)
		switch statusCode {
		case 404:
			return NewJSONError(err, http.StatusNotFound, c)
		default:
			return err
		}
	}

	var editTaskForm EditTaskForm
	err = c.Bind(editTaskForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Update", err, c).Error(err)
		return err
	}
	logging.SharedInstance().MethodInfo("TasksController", "Update", c).Debugf("post update parameter: %+v", editTaskForm)

	valid, err := validators.TaskUpdateValidation(editTaskForm.Title, editTaskForm.Description)
	if err != nil || !valid {
		logging.SharedInstance().MethodInfo("TasksController", "Update", c).Infof("validation error: %v", err)
		return NewJSONError(err, http.StatusUnprocessableEntity, c)
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
		return err
	}

	jsonAllLists, err := allListsResponse(projectService)
	if err != nil {
		return err
	}
	logging.SharedInstance().MethodInfo("TasksController", "Update", c).Debugf("update task: %+v", task)
	logging.SharedInstance().MethodInfo("TasksController", "Update", c).Info("success to update task")
	return c.JSON(http.StatusOK, jsonAllLists)
}

// Delete a task
func (u *Tasks) Delete(c echo.Context) error {
	currentUser, err := LoginRequired(c)
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Delete", c).Infof("loging error: %v", err)
		return NewJSONError(err, http.StatusUnauthorized, c)
	}

	projectService, parentList, statusCode, err := setProjectAndList(c, currentUser)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Delete", err, c).Error(err)
		switch statusCode {
		case 404:
			return NewJSONError(err, http.StatusNotFound, c)
		default:
			return err
		}
	}

	task, statusCode, err := setTask(c, parentList)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("TasksController", "Delete", err, c).Error(err)
		switch statusCode {
		case 404:
			return NewJSONError(err, http.StatusNotFound, c)
		default:
			return err
		}
	}

	err = task.Delete()
	if err != nil {
		logging.SharedInstance().MethodInfo("TasksController", "Delete", c).Info(err)
		return NewJSONError(err, http.StatusBadRequest, c)
	}
	jsonAllLists, err := allListsResponse(projectService)
	if err != nil {
		return err
	}
	logging.SharedInstance().MethodInfo("TasksController", "Delete", c).Info("success to delete a task")
	return c.JSON(http.StatusOK, jsonAllLists)
}

func setProjectAndList(c echo.Context, currentUser *services.User) (*services.Project, *services.List, int, error) {
	projectID, err := strconv.ParseInt(c.Param("project_id"), 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		return nil, nil, 404, err
	}
	projectService, err := handlers.FindProject(projectID)
	if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
		return nil, nil, 404, err
	}
	listID, err := strconv.ParseInt(c.Param("list_id"), 10, 64)
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

func setTask(c echo.Context, list *services.List) (*services.Task, int, error) {
	taskID, err := strconv.ParseInt(c.Param("task_id"), 10, 64)
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

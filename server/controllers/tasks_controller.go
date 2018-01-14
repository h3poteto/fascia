package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/commands/board"
	"github.com/h3poteto/fascia/server/middlewares"
	"github.com/h3poteto/fascia/server/validators"
	"github.com/h3poteto/fascia/server/views"

	"database/sql"
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// Tasks is controller struct for tasks
type Tasks struct {
}

// NewTaskForm is struct for new task
type NewTaskForm struct {
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
}

// MoveTaskForm is struct for move task
type MoveTaskForm struct {
	ToListID     int64 `json:"to_list_id" form:"to_list_id"`
	PrevToTaskID int64 `json:"prev_to_task_id" form:"prev_to_task_id"`
}

// EditTaskForm is struct for edit task
type EditTaskForm struct {
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
}

// Create a new task
func (u *Tasks) Create(c echo.Context) error {
	lc, ok := c.(*middlewares.ListContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	projectService := lc.ProjectService
	parentList := lc.ListService

	newTaskForm := new(NewTaskForm)
	err := c.Bind(newTaskForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Debugf("post new task parameter: %+v", newTaskForm)

	valid, err := validators.TaskCreateValidation(newTaskForm.Title, newTaskForm.Description)
	if err != nil || !valid {
		logging.SharedInstance().Controller(c).Infof("validation error: %v", err)
		return NewValidationError(err, http.StatusUnprocessableEntity, c)
	}

	task := board.NewTask(
		0,
		parentList.ListEntity.ID,
		projectService.ProjectEntity.ID,
		parentList.ListEntity.UserID,
		sql.NullInt64{},
		newTaskForm.Title,
		newTaskForm.Description,
		false,
		sql.NullString{},
	)

	if err := task.Save(); err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Errorf("save failed: %v", err)
		return err
	}

	jsonAllLists, err := allListsResponse(projectService)
	if err != nil {
		return err
	}
	logging.SharedInstance().Controller(c).Debugf("create task success: %+v", task)
	logging.SharedInstance().Controller(c).Info("success to create task")
	return c.JSON(http.StatusOK, jsonAllLists)
}

// Show render json with task detail
func (u *Tasks) Show(c echo.Context) error {
	tc, ok := c.(*middlewares.TaskContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	task := tc.TaskService

	jsonTask, err := views.ParseTaskJSON(task.TaskEntity)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Info("success to get task")
	return c.JSON(http.StatusOK, jsonTask)
}

// MoveTask move a task to another list
func (u *Tasks) MoveTask(c echo.Context) error {
	tc, ok := c.(*middlewares.TaskContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	projectService := tc.ProjectService
	task := tc.TaskService

	moveTaskFrom := new(MoveTaskForm)
	err := c.Bind(moveTaskFrom)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Debugf("post move taks parameter: %+v", moveTaskFrom)

	valid, err := validators.TaskMoveValidation(moveTaskFrom.ToListID, moveTaskFrom.PrevToTaskID)
	if err != nil || !valid {
		logging.SharedInstance().Controller(c).Infof("validation error: %v", err)
		return NewValidationError(err, http.StatusUnprocessableEntity, c)
	}

	var prevToTaskID *int64
	if moveTaskFrom.PrevToTaskID != 0 {
		prevToTaskID = &moveTaskFrom.PrevToTaskID
	}

	if err := task.ChangeList(moveTaskFrom.ToListID, prevToTaskID); err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Errorf("failed change list: %v", err)
		return err
	}

	jsonAllLists, err := allListsResponse(projectService)
	if err != nil {
		return err
	}
	logging.SharedInstance().Controller(c).Debugf("move task: %+v", task)
	logging.SharedInstance().Controller(c).Info("success to move task")
	return c.JSON(http.StatusOK, jsonAllLists)
}

// Update a task
func (u *Tasks) Update(c echo.Context) error {
	tc, ok := c.(*middlewares.TaskContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	projectService := tc.ProjectService
	task := tc.TaskService

	editTaskForm := new(EditTaskForm)
	err := c.Bind(editTaskForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Debugf("post update parameter: %+v", editTaskForm)

	valid, err := validators.TaskUpdateValidation(editTaskForm.Title, editTaskForm.Description)
	if err != nil || !valid {
		logging.SharedInstance().Controller(c).Infof("validation error: %v", err)
		return NewValidationError(err, http.StatusUnprocessableEntity, c)
	}

	err = task.Update(
		task.TaskEntity.ListID,
		task.TaskEntity.IssueNumber,
		editTaskForm.Title,
		editTaskForm.Description,
		task.TaskEntity.PullRequest,
		task.TaskEntity.HTMLURL,
	)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	jsonAllLists, err := allListsResponse(projectService)
	if err != nil {
		return err
	}
	logging.SharedInstance().Controller(c).Debugf("update task: %+v", task)
	logging.SharedInstance().Controller(c).Info("success to update task")
	return c.JSON(http.StatusOK, jsonAllLists)
}

// Delete a task
func (u *Tasks) Delete(c echo.Context) error {
	tc, ok := c.(*middlewares.TaskContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	projectService := tc.ProjectService
	task := tc.TaskService

	err := task.Delete()
	if err != nil {
		logging.SharedInstance().Controller(c).Info(err)
		return NewJSONError(err, http.StatusBadRequest, c)
	}
	jsonAllLists, err := allListsResponse(projectService)
	if err != nil {
		return err
	}
	logging.SharedInstance().Controller(c).Info("success to delete a task")
	return c.JSON(http.StatusOK, jsonAllLists)
}

func allListsResponse(projectService *board.Project) (*views.AllLists, error) {
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

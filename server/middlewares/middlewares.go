package middlewares

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/services"
	"github.com/h3poteto/fascia/server/session"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"

	"fmt"
	"net/http"
	"os"
	"strconv"
)

// JSONError is a struct for http error
type JSONError struct {
	Code    int    `json:code`
	Message string `json:message`
}

// NewJSONError render error json response and return error
func NewJSONError(err error, code int, c echo.Context) error {
	c.JSON(code, &JSONError{
		Code:    code,
		Message: http.StatusText(code),
	})
	return err
}

type LoginContext struct {
	echo.Context
	CurrentUserService *services.User
}

type ProjectContext struct {
	LoginContext
	ProjectService *services.Project
}

type ListContext struct {
	ProjectContext
	ListService *services.List
}

type TaskContext struct {
	ListContext
	TaskService *services.Task
}

func Login() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userService, err := CheckLogin(c)
			if err != nil {
				logging.SharedInstance().Controller(c).Info(err)
				return NewJSONError(err, http.StatusUnauthorized, c)
			}
			uc := &LoginContext{
				c,
				userService,
			}
			return next(uc)
		}
	}
}

// CheckLogin authenticate user
// If unauthorized, return 401
func CheckLogin(c echo.Context) (*services.User, error) {
	id, err := session.SharedInstance().Get(c.Request(), "current_user_id")
	if id == nil {
		return nil, errors.New("not logined")
	}
	currentUser, err := handlers.FindUser(id.(int64))
	if err != nil {
		return nil, err
	}
	return currentUser, nil
}

func Project() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			uc, ok := c.(*LoginContext)
			if !ok {
				err := errors.New("Can not cast context")
				logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
				return err
			}
			projectID, err := strconv.ParseInt(c.Param("project_id"), 10, 64)
			if err != nil {
				err := errors.Wrap(err, "parse error")
				logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
				return NewJSONError(err, http.StatusNotFound, c)
			}
			projectService, err := handlers.FindProject(projectID)
			if err != nil || !(projectService.CheckOwner(uc.CurrentUserService.UserEntity.UserModel.ID)) {
				logging.SharedInstance().Controller(c).Warnf("project not found: %v", err)
				return NewJSONError(err, http.StatusNotFound, c)
			}

			pc := &ProjectContext{
				*uc,
				projectService,
			}
			return next(pc)
		}
	}
}

func List() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			pc, ok := c.(*ProjectContext)
			if !ok {
				err := errors.New("Can not cast context")
				logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
				return err
			}
			listID, err := strconv.ParseInt(c.Param("list_id"), 10, 64)
			if err != nil {
				err := errors.Wrap(err, "parse error")
				logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
				return NewJSONError(err, http.StatusNotFound, c)
			}
			listService, err := handlers.FindList(pc.ProjectService.ProjectEntity.ProjectModel.ID, listID)
			if err != nil {
				logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
				return NewJSONError(err, http.StatusNotFound, c)
			}
			lc := &ListContext{
				*pc,
				listService,
			}
			return next(lc)
		}
	}
}

func Task() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			lc, ok := c.(*ListContext)
			if !ok {
				err := errors.New("Can not cast context")
				logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
				return err
			}
			taskID, err := strconv.ParseInt(c.Param("task_id"), 10, 64)
			if err != nil {
				err := errors.Wrap(err, "parse error")
				logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
				return NewJSONError(err, http.StatusNotFound, c)
			}
			taskService, err := handlers.FindTask(lc.ListService.ListEntity.ListModel.ID, taskID)
			if err != nil {
				logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
				return NewJSONError(err, http.StatusNotFound, c)
			}
			tc := &TaskContext{
				*lc,
				taskService,
			}
			return next(tc)
		}
	}
}

// PanicRecover prepare original panic recover using logrus
func PanicRecover() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					var err error
					switch r := r.(type) {
					case error:
						err = r
					default:
						err = errors.Errorf("%v", r)
					}
					logging.SharedInstance().PanicRecover(c).Error(err)
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}

func CustomizeLogger() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: printColored("status") + "=${status} " + printColored("method") + "=${method} " + printColored("path") + "=${uri} " + printColored("requestID") + "=${id} " + printColored("latency") + "=${latency_human} " + printColored("time") + "=${time_rfc3339_nano}\n",
		Output: os.Stdout,
	})
}

func printColored(str string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", 34, str)
}

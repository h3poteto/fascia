package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/domains/list"
	"github.com/h3poteto/fascia/server/domains/project"
	"github.com/h3poteto/fascia/server/domains/task"
	"github.com/h3poteto/fascia/server/domains/user"
	"github.com/h3poteto/fascia/server/session"
	"github.com/h3poteto/fascia/server/usecases/account"
	"github.com/h3poteto/fascia/server/usecases/board"
	"github.com/h3poteto/fascia/server/validators"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"

	"fmt"
	"net/http"
	"os"
	"strconv"
)

// JSONError is a struct for http error
type JSONError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewJSONError render error json response and return error
func NewJSONError(err error, code int, c echo.Context) error {
	c.JSON(code, &JSONError{
		Code:    code,
		Message: http.StatusText(code),
	})
	return err
}

// NewValidationError is as struct for validation error with json
func NewValidationError(err error, code int, c echo.Context) error {
	c.JSON(code, validators.ErrorsByField(err))
	return err
}

// LoginContext prepare login information for users
type LoginContext struct {
	echo.Context
	CurrentUser *user.User
}

// ProjectContext prepare a project service
type ProjectContext struct {
	LoginContext
	Project *project.Project
}

// ListContext prepare a list service
type ListContext struct {
	ProjectContext
	List *list.List
}

// TaskContext prepare a task service
type TaskContext struct {
	ListContext
	Task *task.Task
}

// Login requires login session
// If unauthorized, return 401
func Login() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, err := CheckLogin(c)
			if err != nil {
				logging.SharedInstance().Controller(c).Info(err)
				return NewJSONError(err, http.StatusUnauthorized, c)
			}
			uc := &LoginContext{
				c,
				user,
			}
			return next(uc)
		}
	}
}

// CheckLogin authenticate user
func CheckLogin(c echo.Context) (*user.User, error) {
	id, err := checkJWT(c)
	if err != nil || id == -1 {
		id, err = checkSession(c)
		if err != nil || id == -1 {
			return nil, errors.New("not logged in")
		}
	}

	currentUser, err := account.FindUser(id)
	if err != nil {
		return nil, err
	}
	return currentUser, nil
}

func checkJWT(c echo.Context) (int64, error) {
	u := c.Get("user")
	if u == nil {
		return -1, errors.New("JWT does not exist")
	}
	ut := u.(*jwt.Token)
	claims := ut.Claims.(*config.JwtCustomClaims)
	return claims.CurrentUserID, nil
}

func checkSession(c echo.Context) (int64, error) {
	session_id, err := session.SharedInstance().Get(c.Request(), "current_user_id")
	if err != nil || session_id == nil {
		return -1, errors.New("Session does not exist")
	}
	return session_id.(int64), nil
}

// Project requires a project from project_id
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
			p, err := board.FindProject(projectID)
			if err != nil || !(p.CheckOwner(uc.CurrentUser.ID)) {
				logging.SharedInstance().Controller(c).Warnf("project not found: %v", err)
				return NewJSONError(err, http.StatusNotFound, c)
			}

			pc := &ProjectContext{
				*uc,
				p,
			}
			return next(pc)
		}
	}
}

// List require a list from list_id
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
			l, err := board.FindList(pc.Project.ID, listID)
			if err != nil {
				logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
				return NewJSONError(err, http.StatusNotFound, c)
			}
			lc := &ListContext{
				*pc,
				l,
			}
			return next(lc)
		}
	}
}

// Task requires a task from task_id
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
			t, err := board.FindTask(taskID)
			if err != nil {
				logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
				return NewJSONError(err, http.StatusNotFound, c)
			}
			tc := &TaskContext{
				*lc,
				t,
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

// CustomizeLogger prepqre my logger for echo
func CustomizeLogger() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: printColored("status") + "=${status} " + printColored("method") + "=${method} " + printColored("path") + "=${uri} " + printColored("requestID") + "=${id} " + printColored("latency") + "=${latency_human} " + printColored("time") + "=${time_rfc3339_nano}\n",
		Output: os.Stdout,
	})
}

func printColored(str string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", 34, str)
}

type fundamental interface {
	StackTrace() errors.StackTrace
}

// ErrorLogging logging error and call default error handler in echo
func ErrorLogging(e *echo.Echo) func(error, echo.Context) {
	return func(err error, c echo.Context) {
		// pkg/errorsにより生成されたエラーについては，各コントローラで適切にハンドリングすること
		// ここでは予定外のエラーが発生した場合にログを飛ばしたい
		// 予定外のエラーなので，errors.fundamentalとecho.HTTPError以外のエラーだけを拾えれば十分なはずである
		_, isFundamental := err.(fundamental)
		_, isHTTPError := err.(*echo.HTTPError)
		if !isFundamental && !isHTTPError {
			logging.SharedInstance().Controller(c).Error(err)
		}
		e.DefaultHTTPErrorHandler(err, c)
	}
}

// JWTSkipper skip jwt middleware when auth session exists.
func JWTSkipper(c echo.Context) bool {
	_, err := checkSession(c)
	if err != nil {
		return false
	}
	return true
}

package middlewares

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	_ "github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/services"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"

	"fmt"
	"os"
)

type ProjectContext struct {
	echo.Context
	ProjectService *services.Project
}

type ListContext struct {
	echo.Context
	ProjectService *services.Project
	ListService    *services.List
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

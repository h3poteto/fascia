package logging

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/heroku/rollrus"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// LogStruct provides logger object
type LogStruct struct {
	Log *logrus.Logger
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

var sharedInstance *LogStruct = New()

// New returns a LogStruct
func New() *LogStruct {
	goenv := os.Getenv("APPENV")
	log := logrus.New()
	log.Out = os.Stdout
	if goenv == "production" {
		log.Level = logrus.InfoLevel
	} else {
		log.Level = logrus.DebugLevel
	}
	hook := rollrus.NewHook(os.Getenv("ROLLBAR_TOKEN"), goenv)
	log.Hooks.Add(hook)
	return &LogStruct{Log: log}
}

// SharedInstance returns a singleton object
func SharedInstance() *LogStruct {
	return sharedInstance
}

// MethodInfo is prepare logrus entry with fields
func (u *LogStruct) MethodInfo(model string, action string) *logrus.Entry {
	return u.Log.WithFields(logrus.Fields{
		"time":   time.Now(),
		"model":  model,
		"action": action,
	})
}

// MethodInfoWithStacktrace is prepare logrus entry with fields
func (u *LogStruct) MethodInfoWithStacktrace(model string, action string, err error) *logrus.Entry {
	stackErr, ok := err.(stackTracer)
	if !ok {
		panic("oops, err does not implement Stacktrace")
	}
	st := stackErr.StackTrace()
	traceLength := len(st)
	if traceLength > 5 {
		traceLength = 5
	}

	return u.Log.WithFields(logrus.Fields{
		"time":       time.Now(),
		"model":      model,
		"action":     action,
		"stacktrace": fmt.Sprintf("%+v", st[0:traceLength]),
	})
}

// PanicRecover send error and stacktrace
func (u *LogStruct) PanicRecover(context echo.Context) *logrus.Entry {
	requestID := context.Response().Header().Get(echo.HeaderXRequestID)
	userAgent := context.Request().Header.Get("User-Agent")
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, false)
	return u.Log.WithFields(logrus.Fields{
		"time":       time.Now(),
		"requestID":  requestID,
		"User-Agent": userAgent,
		"model":      "main",
		"stacktrace": string(buf),
	})
}

// Controller is prepare logrus entry with fields
func (u *LogStruct) Controller(context echo.Context) *logrus.Entry {
	requestID := context.Response().Header().Get(echo.HeaderXRequestID)
	userAgent := context.Request().Header.Get("User-Agent")
	return u.Log.WithFields(logrus.Fields{
		"time":       time.Now(),
		"method":     context.Request().Method,
		"requestID":  requestID,
		"User-Agent": userAgent,
		"path":       context.Path(),
	})
}

// ControllerWithStacktrace is prepare logrus entry with fields
func (u *LogStruct) ControllerWithStacktrace(err error, context echo.Context) *logrus.Entry {
	requestID := context.Response().Header().Get(echo.HeaderXRequestID)
	userAgent := context.Request().Header.Get("User-Agent")

	stackErr, ok := err.(stackTracer)
	if !ok {
		return u.Log.WithFields(logrus.Fields{
			"time":       time.Now(),
			"method":     context.Request().Method,
			"requestID":  requestID,
			"User-Agent": userAgent,
			"path":       context.Path(),
			"stacktrace": "oops, err does not implement Stacktrace",
		})
	}
	st := stackErr.StackTrace()
	traceLength := len(st)
	if traceLength > 5 {
		traceLength = 5
	}

	return u.Log.WithFields(logrus.Fields{
		"time":       time.Now(),
		"method":     context.Request().Method,
		"requestID":  requestID,
		"User-Agent": userAgent,
		"path":       context.Path(),
		"stacktrace": fmt.Sprintf("%+v", st[0:traceLength]),
	})
}

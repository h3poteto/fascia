package logging

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/johntdyer/slackrus"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type LogStruct struct {
	Log *logrus.Logger
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

var sharedInstance *LogStruct = New()

func New() *LogStruct {
	goenv := os.Getenv("APPENV")
	log := logrus.New()
	log.Out = os.Stdout
	if goenv == "production" {
		log.Level = logrus.InfoLevel
		log.Hooks.Add(&slackrus.SlackrusHook{
			HookURL:        os.Getenv("SLACK_URL"),
			AcceptedLevels: slackrus.LevelThreshold(logrus.ErrorLevel),
			Channel:        "#fascia",
			IconEmoji:      ":bapho:",
			Username:       "logrus",
		})
	} else {
		log.Level = logrus.DebugLevel
	}
	return &LogStruct{Log: log}
}

func SharedInstance() *LogStruct {
	return sharedInstance
}

// MethodInfo is prepare logrus entry with fields
func (u *LogStruct) MethodInfo(model string, action string, _ ...echo.Context) *logrus.Entry {
	return u.Log.WithFields(logrus.Fields{
		"time":   time.Now(),
		"model":  model,
		"action": action,
	})
}

// MethodInfoWithStacktrace is prepare logrus entry with fields
func (u *LogStruct) MethodInfoWithStacktrace(model string, action string, err error, _ ...echo.Context) *logrus.Entry {
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
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, false)
	return u.Log.WithFields(logrus.Fields{
		"time":       time.Now(),
		"requestID":  requestID,
		"model":      "main",
		"stacktrace": string(buf),
	})
}

func (u *LogStruct) Controller(context echo.Context) *logrus.Entry {
	requestID := context.Response().Header().Get(echo.HeaderXRequestID)

	return u.Log.WithFields(logrus.Fields{
		"time":      time.Now(),
		"requestID": requestID,
		"path":      context.Path(),
	})
}

func (u *LogStruct) ControllerWithStacktrace(err error, context echo.Context) *logrus.Entry {
	requestID := context.Response().Header().Get(echo.HeaderXRequestID)

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
		"requestID":  requestID,
		"path":       context.Path(),
		"stacktrace": fmt.Sprintf("%+v", st[0:traceLength]),
	})
}

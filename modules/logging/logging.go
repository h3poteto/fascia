package logging

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/johntdyer/slackrus"
	"github.com/pkg/errors"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

type LogStruct struct {
	Log *logrus.Logger
}

type Stacktrace interface {
	Stacktrace() []errors.Frame
}

var sharedInstance *LogStruct = New()

func New() *LogStruct {
	goenv := os.Getenv("GOJIENV")
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
func (u *LogStruct) MethodInfo(model string, action string, stack bool, context ...web.C) *logrus.Entry {
	requestID := "null"
	if len(context) > 0 {
		requestID = middleware.GetReqID(context[0])
	}
	if stack {
		buf := make([]byte, 1024)
		runtime.Stack(buf, false)
		return u.Log.WithFields(logrus.Fields{
			"time":       time.Now(),
			"requestID":  requestID,
			"model":      model,
			"action":     action,
			"stacktrace": string(buf),
		})
	}
	return u.Log.WithFields(logrus.Fields{
		"time":      time.Now(),
		"requestID": requestID,
		"model":     model,
		"action":    action,
	})
}

func (u *LogStruct) MethodInfoWithStacktrace(model string, action string, err error, context ...web.C) *logrus.Entry {
	requestID := "null"
	if len(context) > 0 {
		requestID = middleware.GetReqID(context[0])
	}

	stackErr, ok := err.(Stacktrace)
	if !ok {
		panic("oops, err does not implement Stacktrace")
	}
	st := stackErr.Stacktrace()

	return u.Log.WithFields(logrus.Fields{
		"time":       time.Now(),
		"requestID":  requestID,
		"model":      model,
		"action":     action,
		"stacktrace": fmt.Sprintf("%+v", st[0:5]),
	})
}

// PanicRecover send error and stacktrace
func (u *LogStruct) PanicRecover(context web.C) *logrus.Entry {
	requestID := middleware.GetReqID(context)
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, false)
	return u.Log.WithFields(logrus.Fields{
		"time":       time.Now(),
		"requestID":  requestID,
		"model":      "main",
		"stacktrace": string(buf),
	})
}

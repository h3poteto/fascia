package logging

import (
	"github.com/Sirupsen/logrus"
	"github.com/johntdyer/slackrus"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"os"
	"runtime"
	"time"
)

type LogStruct struct {
	log *logrus.Logger
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
	return &LogStruct{log: log}
}

func SharedInstance() *LogStruct {
	return sharedInstance
}

// MethodInfo is prepare logrus entry with fields
func (u *LogStruct) MethodInfo(model string, method string, stack bool, context ...web.C) *logrus.Entry {
	requestID := "null"
	if len(context) > 0 {
		requestID = middleware.GetReqID(context[0])
	}
	if stack {
		buf := make([]byte, 1024)
		runtime.Stack(buf, false)
		return u.log.WithFields(logrus.Fields{
			"time":       time.Now(),
			"requestID":  requestID,
			"model":      model,
			"method":     method,
			"stacktrace": string(buf),
		})
	}
	return u.log.WithFields(logrus.Fields{
		"time":      time.Now(),
		"requestID": requestID,
		"model":     model,
		"method":    method,
	})
}

// PanicRecover send error and stacktrace
func (u *LogStruct) PanicRecover(context web.C) *logrus.Entry {
	requestID := middleware.GetReqID(context)
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, false)
	return u.log.WithFields(logrus.Fields{
		"time":       time.Now(),
		"requestID":  requestID,
		"model":      "main",
		"stacktrace": string(buf),
	})
}

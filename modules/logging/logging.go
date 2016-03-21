package logging

import (
	"github.com/Sirupsen/logrus"
	"github.com/johntdyer/slackrus"
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

func (u *LogStruct) MethodInfo(model string, method string, stack ...bool) *logrus.Entry {
	if len(stack) > 0 && stack[0] {
		buf := make([]byte, 1<<16)
		runtime.Stack(buf, false)
		return u.log.WithFields(logrus.Fields{
			"time":       time.Now(),
			"model":      model,
			"method":     method,
			"stacktrace": string(buf),
		})
	}
	return u.log.WithFields(logrus.Fields{
		"time":   time.Now(),
		"model":  model,
		"method": method,
	})
}

func (u *LogStruct) PanicRecover() *logrus.Entry {
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, false)
	return u.log.WithFields(logrus.Fields{
		"time":       time.Now(),
		"model":      "main",
		"stacktrace": string(buf),
	})
}

package logging

import (
	"github.com/Sirupsen/logrus"
	"github.com/johntdyer/slackrus"
	"os"
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
	} else {
		log.Level = logrus.DebugLevel
	}
	log.Hooks.Add(&slackrus.SlackrusHook{
		HookURL:        os.Getenv("SLACK_URL"),
		AcceptedLevels: slackrus.LevelThreshold(logrus.PanicLevel),
		Channel:        "#fascia",
		IconEmoji:      ":bapho:",
		Username:       "logrus",
	})

	return &LogStruct{log: log}
}

func SharedInstance() *LogStruct {
	return sharedInstance
}

func (u *LogStruct) MethodInfo(model string, method string) *logrus.Entry {
	return u.log.WithFields(logrus.Fields{
		"model":  model,
		"method": method,
	})
}

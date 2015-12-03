package logging

import (
	"github.com/Sirupsen/logrus"
	"os"
)

type LogStruct struct {
	log *logrus.Logger
}

var sharedInstance *LogStruct = New()

func New() *LogStruct {
	log := logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.DebugLevel
	return &LogStruct{log: log}
}

func SharedInstance() *LogStruct {
	return sharedInstance
}

func (u *LogStruct) BaseInfo(model string, method string) *logrus.Entry {
	return u.log.WithFields(logrus.Fields{
		"model":  model,
		"method": method,
	})
}

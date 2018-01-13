package handlers

import (
	"github.com/h3poteto/fascia/server/commands/project"
)

// FindTask returns a task service
func FindTask(listID, taskID int64) (*project.Task, error) {
	return project.FindTask(listID, taskID)
}

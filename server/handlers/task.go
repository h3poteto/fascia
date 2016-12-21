package handlers

import (
	"github.com/h3poteto/fascia/server/services"
)

// FindTask returns a task service
func FindTask(listID, taskID int64) (*services.Task, error) {
	return services.FindTask(listID, taskID)
}

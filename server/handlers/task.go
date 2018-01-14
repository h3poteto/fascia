package handlers

import (
	"github.com/h3poteto/fascia/server/commands/board"
)

// FindTask returns a task service
func FindTask(listID, taskID int64) (*board.Task, error) {
	return board.FindTask(listID, taskID)
}

package handlers

import (
	"database/sql"

	"github.com/h3poteto/fascia/server/commands/board"
)

// NewList returns a new list service
func NewList(id, projectID, userID int64, title, color string, optionID sql.NullInt64, isHidden bool) *board.List {
	return board.NewList(id, projectID, userID, title, color, optionID, isHidden)
}

// FindList returns a list service
func FindList(projectID, listID int64) (*board.List, error) {
	return board.FindListByID(projectID, listID)
}

package handlers

import (
	"database/sql"

	"github.com/h3poteto/fascia/server/commands/project"
)

// NewList returns a new list service
func NewList(id, projectID, userID int64, title, color string, optionID sql.NullInt64, isHidden bool) *project.List {
	return project.NewList(id, projectID, userID, title, color, optionID, isHidden)
}

// FindList returns a list service
func FindList(projectID, listID int64) (*project.List, error) {
	return project.FindListByID(projectID, listID)
}

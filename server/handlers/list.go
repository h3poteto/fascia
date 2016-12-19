package handlers

import (
	"database/sql"

	"github.com/h3poteto/fascia/server/services"
)

func NewList(id, projectID, userID int64, title, color string, optionID sql.NullInt64, isHidden bool) *services.List {
	return services.NewList(id, projectID, userID, title, color, optionID, isHidden)
}

func FindList(projectID, listID int64) (*services.List, error) {
	return services.FindListByID(projectID, listID)
}

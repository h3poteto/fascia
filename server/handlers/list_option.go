package handlers

import (
	"github.com/h3poteto/fascia/server/commands/project"
)

// ListOptionAll returns all list options
func ListOptionAll() ([]*project.ListOption, error) {
	return project.ListOptionAll()
}

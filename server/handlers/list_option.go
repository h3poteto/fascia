package handlers

import (
	"github.com/h3poteto/fascia/server/services"
)

// ListOptionAll returns all list options
func ListOptionAll() ([]*services.ListOption, error) {
	return services.ListOptionAll()
}

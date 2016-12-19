package handlers

import (
	"github.com/h3poteto/fascia/server/services"
)

func ListOptionAll() ([]*services.ListOption, error) {
	return services.ListOptionAll()
}

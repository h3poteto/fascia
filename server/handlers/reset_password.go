package handlers

import (
	"github.com/h3poteto/fascia/server/services"
)

func ChangeUserPassword(id int64, token string, password string) (*services.User, error) {
	return services.ChangeUserPassword(id, token, password)
}

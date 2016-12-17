package handlers

import (
	"github.com/h3poteto/fascia/server/services"
)

func FindUser(id int64) (*services.User, error) {
	return services.FindUser(id)
}

func LoginUser(email string, password string) (*services.User, error) {
	return services.LoginUser(email, password)
}

func FindOrCreateUserFromGithub(token string) (*services.User, error) {
	return services.FindOrCreateUserFromGithub(token)
}

package handlers

import (
	"github.com/h3poteto/fascia/server/services"
)

func RegistrationUser(email, password, passwordConfirm string) (*services.User, error) {
	return services.RegistrationUser(email, password, passwordConfirm)
}

func FindUser(id int64) (*services.User, error) {
	return services.FindUser(id)
}

func FindUserByEmail(email string) (*services.User, error) {
	return services.FindUserByEmail(email)
}

func LoginUser(email string, password string) (*services.User, error) {
	return services.LoginUser(email, password)
}

func FindOrCreateUserFromGithub(token string) (*services.User, error) {
	return services.FindOrCreateUserFromGithub(token)
}

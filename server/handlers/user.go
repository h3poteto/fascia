package handlers

import (
	"github.com/h3poteto/fascia/server/services"
)

// RegistrationUser create a user with email and password
func RegistrationUser(email, password, passwordConfirm string) (*services.User, error) {
	return services.RegistrationUser(email, password, passwordConfirm)
}

// FindUser search a user service
func FindUser(id int64) (*services.User, error) {
	return services.FindUser(id)
}

// FindUserByEmail search a user service according to email
func FindUserByEmail(email string) (*services.User, error) {
	return services.FindUserByEmail(email)
}

// LoginUser authenticate with email and password
func LoginUser(email string, password string) (*services.User, error) {
	return services.LoginUser(email, password)
}

// FindOrCreateUserFromGithub authenticate with github, and return a user service
func FindOrCreateUserFromGithub(token string) (*services.User, error) {
	return services.FindOrCreateUserFromGithub(token)
}

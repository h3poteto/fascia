package handlers

import (
	"github.com/h3poteto/fascia/server/commands/account"
)

// RegistrationUser create a user with email and password
func RegistrationUser(email, password, passwordConfirm string) (*account.User, error) {
	return account.RegistrationUser(email, password, passwordConfirm)
}

// FindUser search a user service
func FindUser(id int64) (*account.User, error) {
	return account.FindUser(id)
}

// FindUserByEmail search a user service according to email
func FindUserByEmail(email string) (*account.User, error) {
	return account.FindUserByEmail(email)
}

// LoginUser authenticate with email and password
func LoginUser(email string, password string) (*account.User, error) {
	return account.LoginUser(email, password)
}

// FindOrCreateUserFromGithub authenticate with github, and return a user service
func FindOrCreateUserFromGithub(token string) (*account.User, error) {
	return account.FindOrCreateUserFromGithub(token)
}

package handlers

import (
	"github.com/h3poteto/fascia/server/commands/account"
)

// ChangeUserPassword change owner user password
func ChangeUserPassword(id int64, token string, password string) (*account.User, error) {
	return account.ChangeUserPassword(id, token, password)
}

// GenerateResetPassword returns a reset password service
func GenerateResetPassword(userID int64, email string) (*account.ResetPassword, error) {
	return account.GenerateResetPassword(userID, email)
}

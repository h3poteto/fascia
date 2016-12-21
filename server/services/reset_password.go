package services

import (
	"github.com/h3poteto/fascia/server/entities/reset_password"
)

// ResetPassword has a reset password entity
type ResetPassword struct {
	ResetPasswordEntity *reset_password.ResetPassword
}

// ChangeUserPassword
func ChangeUserPassword(id int64, token string, password string) (*User, error) {
	// reset_passwordモデルを探し出す必要がある
	r, err := reset_password.FindAvailable(id, token)
	if err != nil {
		return nil, err
	}

	u, err := r.ChangeUserPassword(password)
	if err != nil {
		return nil, err
	}
	return NewUser(u), nil
}

// GenerateResetPassword create new token and returns a reset password service
func GenerateResetPassword(userID int64, email string) (*ResetPassword, error) {
	r, err := reset_password.GenerateResetPassword(userID, email)
	if err != nil {
		return nil, err
	}
	return &ResetPassword{
		ResetPasswordEntity: r,
	}, nil
}

// AuthenticateResetPassword check token
func AuthenticateResetPassword(id int64, token string) error {
	return reset_password.Authenticate(id, token)
}

// Save save a reset password entity
func (r *ResetPassword) Save() error {
	return r.ResetPasswordEntity.Save()
}

package services

import (
	"github.com/h3poteto/fascia/server/aggregations/reset_password"
)

type ResetPassword struct {
	ResetPasswordAggregation *reset_password.ResetPassword
}

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
	return NewUserService(u), nil
}

func GenerateResetPassword(userID int64, email string) (*ResetPassword, error) {
	r, err := reset_password.GenerateResetPassword(userID, email)
	if err != nil {
		return nil, err
	}
	return &ResetPassword{
		ResetPasswordAggregation: r,
	}, nil
}

func AuthenticateResetPassword(id int64, token string) error {
	return reset_password.Authenticate(id, token)
}

func (r *ResetPassword) Save() error {
	return r.ResetPasswordAggregation.Save()
}

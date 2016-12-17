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

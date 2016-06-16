package validators

import (
	"github.com/pkg/errors"
)

type passwordCreate struct {
	Email string `valid:"email"`
}

type passwordUpdate struct {
	ResetToken      string `valid:"required"`
	Password        string `valid:"min=6,max=255"`
	PasswordConfirm string `valid:"min=6,max=255"`
}

func PasswordCreateValidation(email string) error {
	form := &passwordCreate{
		Email: email,
	}
	return validate.Struct(form)
}

func PasswordUpdateValidation(resetToken string, password string, passwordConfirm string) error {
	if password != passwordConfirm {
		return errors.New("password and password confirm did not match")
	}
	form := &passwordUpdate{
		ResetToken:      resetToken,
		Password:        password,
		PasswordConfirm: passwordConfirm,
	}
	return validate.Struct(form)
}

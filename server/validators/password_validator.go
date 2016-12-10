package validators

import (
	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
)

type passwordCreate struct {
	Email string `valid:"email,required"`
}

type passwordUpdate struct {
	ResetToken      string `valid:"required"`
	Password        string `valid:"stringlength(6|255)"`
	PasswordConfirm string `valid:"stringlength(6|255)"`
}

// PasswordCreateValidation check form variable when create reset_passwords
func PasswordCreateValidation(email string) (bool, error) {
	form := &passwordCreate{
		Email: email,
	}
	return govalidator.ValidateStruct(form)
}

// PasswordUpdateValidation check form variable when update reset_passwords
func PasswordUpdateValidation(resetToken string, password string, passwordConfirm string) (bool, error) {
	if password != passwordConfirm {
		return false, errors.New("password and password confirm did not match")
	}
	form := &passwordUpdate{
		ResetToken:      resetToken,
		Password:        password,
		PasswordConfirm: passwordConfirm,
	}
	return govalidator.ValidateStruct(form)
}

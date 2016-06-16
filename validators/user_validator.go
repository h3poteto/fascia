package validators

import (
	"github.com/pkg/errors"
)

type userRegistration struct {
	Email           string `valid:"email"`
	Password        string `valid:"min=6,max=255"`
	PasswordConfirm string `valid:"min=6,max=255"`
}

func UserRegistrationValidation(email string, password string, passwordConfirm string) error {
	if password != passwordConfirm {
		return errors.New("password and password confirm did not match")
	}
	form := &userRegistration{
		Email:           email,
		Password:        password,
		PasswordConfirm: passwordConfirm,
	}
	return validate.Struct(form)
}

package validators

import (
	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
)

type userRegistration struct {
	Email           string `valid:"email,required"`
	Password        string `valid:"length(4|255)"`
	PasswordConfirm string `valid:"length(4|255)"`
}

func UserRegistrationValidation(email string, password string, passwordConfirm string) (bool, error) {
	if password != passwordConfirm {
		return false, errors.New("password and password confirm did not match")
	}
	form := &userRegistration{
		Email:           email,
		Password:        password,
		PasswordConfirm: passwordConfirm,
	}
	return govalidator.ValidateStruct(form)
}

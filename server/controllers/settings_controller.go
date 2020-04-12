package controllers

import (
	"net/http"

	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/middlewares"
	"github.com/h3poteto/fascia/server/usecases/account"
	"github.com/h3poteto/fascia/server/views"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Settings type of settings controller
type Settings struct{}

// PasswordForm accepts update password requests
type PasswordForm struct {
	Password        string `json:"password" form:"password"`
	PasswordConfirm string `json:"password_confirm" form:"password_confirm"`
}

// Password updates password for the user.
func (s *Settings) Password(c echo.Context) error {
	uc, ok := c.(*middlewares.LoginContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	currentUser := uc.CurrentUser

	updatePasswordForm := new(PasswordForm)
	err := c.Bind(updatePasswordForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	u, err := account.UpdatePassword(currentUser.ID, updatePasswordForm.Password, updatePasswordForm.PasswordConfirm)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	jsonUser, err := views.ParseUserJSON(u)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	return c.JSON(http.StatusOK, jsonUser)
}

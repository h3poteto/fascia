package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/commands/account"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/mailers/password_mailer"
	"github.com/h3poteto/fascia/server/validators"

	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// Passwords is controller struct for passwords
type Passwords struct {
}

// NewPasswordForm is struct for new password
type NewPasswordForm struct {
	Email string `form:"email"`
	Token string `form:"token"`
}

// EditPasswordForm is struct for edit password
type EditPasswordForm struct {
	Token           string `form:"token"`
	ResetToken      string `form:"reset_token"`
	Password        string `form:"password"`
	PasswordConfirm string `form:"password_confirm"`
}

// New renders a new password form
func (u *Passwords) New(c echo.Context) error {
	token, err := GenerateCSRFToken(c)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	return c.Render(http.StatusOK, "new_password.html.tpl", map[string]interface{}{
		"title": "PasswordReset",
		"token": token,
	})
}

// Create a new password
func (u *Passwords) Create(c echo.Context) error {
	newPasswordForm := new(NewPasswordForm)
	err := c.Bind(newPasswordForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	if !CheckCSRFToken(c, newPasswordForm.Token) {
		err := errors.New("cannot verify CSRF token")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	valid, err := validators.PasswordCreateValidation(newPasswordForm.Email)
	if err != nil || !valid {
		logging.SharedInstance().Controller(c).Infof("validation failed: %v", err)
		return c.Redirect(http.StatusFound, "/passwords/new")
	}

	targetUser, err := handlers.FindUserByEmail(newPasswordForm.Email)
	if err != nil {
		// OKにしておかないとEmail探りに使われる
		logging.SharedInstance().Controller(c).Infof("cannot find user: %v", err)
		return c.Redirect(http.StatusFound, "/sign_in")
	}

	reset, err := handlers.GenerateResetPassword(targetUser.UserEntity.ID, targetUser.UserEntity.Email)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	if err := reset.Save(); err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	// ここでemail送信
	go password_mailer.Reset(reset.ResetPasswordEntity.ID, targetUser.UserEntity.Email, reset.ResetPasswordEntity.Token)
	logging.SharedInstance().Controller(c).Info("success to send password reset request")
	return c.Redirect(http.StatusFound, "/sign_in")
}

// Edit renders a edit password form
func (u *Passwords) Edit(c echo.Context) error {
	token, err := GenerateCSRFToken(c)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	resetToken := c.QueryParam("token")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return NewJSONError(err, http.StatusNotFound, c)
	}
	if err := account.AuthenticateResetPassword(id, resetToken); err != nil {
		logging.SharedInstance().Controller(c).Infof("cannot authenticate reset password: %v", err)
		return err
	}
	return c.Render(http.StatusOK, "edit_password.html.tpl", map[string]interface{}{
		"title":      "PasswordReset",
		"token":      token,
		"id":         id,
		"resetToken": resetToken,
	})
}

// Update a password
func (u *Passwords) Update(c echo.Context) error {
	editPasswordForm := new(EditPasswordForm)
	err := c.Bind(editPasswordForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameters")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	if !CheckCSRFToken(c, editPasswordForm.Token) {
		err := errors.New("cannot verify CSRF token")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return NewJSONError(err, http.StatusNotFound, c)
	}

	valid, err := validators.PasswordUpdateValidation(editPasswordForm.ResetToken, editPasswordForm.Password, editPasswordForm.PasswordConfirm)
	if err != nil || !valid {
		logging.SharedInstance().Controller(c).Infof("validation failed: %v", err)
		return c.Redirect(http.StatusFound, "/passwords/"+string(id)+"/edit")
	}

	targetUser, err := handlers.ChangeUserPassword(id, editPasswordForm.ResetToken, editPasswordForm.Password)
	if err != nil {
		logging.SharedInstance().Controller(c).Infof("cannot authenticate reset password: %v", err)
		return err
	}

	go password_mailer.Changed(targetUser.UserEntity.Email)
	logging.SharedInstance().Controller(c).Info("success to change password")
	return c.Redirect(http.StatusFound, "/sign_in")
}

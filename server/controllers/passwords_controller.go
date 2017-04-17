package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/mailers/password_mailer"
	"github.com/h3poteto/fascia/server/services"
	"github.com/h3poteto/fascia/server/validators"

	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type Passwords struct {
}

type NewPasswordForm struct {
	Email string `param:"email"`
	Token string `param:"token"`
}

type EditPasswordForm struct {
	Token           string `param:"token"`
	ResetToken      string `param:"reset_token"`
	Password        string `param:"password"`
	PasswordConfirm string `param:"password_confirm"`
}

// tokenを発行し，expireと合わせてreset_passwordモデルにDB保存する
// idとtokenをメールで送る
// idとtoken, expireがあっていたらpasswordの編集を許可する
// passwordを新たに保存する
func (u *Passwords) New(c echo.Context) error {
	token, err := GenerateCSRFToken(c)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "New", err, c).Error(err)
		return err
	}
	return c.Render(http.StatusOK, "new_password.html.tpl", map[string]interface{}{
		"title": "PasswordReset",
		"token": token,
	})
}

func (u *Passwords) Create(c echo.Context) error {
	var newPasswordForm NewPasswordForm
	err := c.Bind(newPasswordForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "Create", err, c).Error(err)
		return err
	}

	if !CheckCSRFToken(c, newPasswordForm.Token) {
		err := errors.New("cannot verify CSRF token")
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "Create", err, c).Error(err)
		return err
	}

	valid, err := validators.PasswordCreateValidation(newPasswordForm.Email)
	if err != nil || !valid {
		logging.SharedInstance().MethodInfo("PasswordsController", "Create", c).Infof("validation failed: %v", err)
		return c.Redirect(http.StatusFound, "/passwords/new")
	}

	targetUser, err := handlers.FindUserByEmail(newPasswordForm.Email)
	if err != nil {
		// OKにしておかないとEmail探りに使われる
		logging.SharedInstance().MethodInfo("PasswordsController", "Create", c).Infof("cannot find user: %v", err)
		return c.Redirect(http.StatusFound, "/sign_in")
	}

	reset, err := handlers.GenerateResetPassword(targetUser.UserEntity.UserModel.ID, targetUser.UserEntity.UserModel.Email)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "Create", err, c).Error(err)
		return err
	}
	if err := reset.Save(); err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "Create", err, c).Error(err)
		return err
	}
	// ここでemail送信
	go password_mailer.Reset(reset.ResetPasswordEntity.ResetPasswordModel.ID, targetUser.UserEntity.UserModel.Email, reset.ResetPasswordEntity.ResetPasswordModel.Token)
	logging.SharedInstance().MethodInfo("PasswordsController", "Create", c).Info("success to send password reset request")
	return c.Redirect(http.StatusFound, "/sign_in")
}

func (u *Passwords) Edit(c echo.Context) error {
	token, err := GenerateCSRFToken(c)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "Edit", err, c).Error(err)
		return err
	}
	resetToken := c.QueryParam("token")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "Edit", err, c).Error(err)
		return c.JSON(http.StatusNotFound, &JSONError{message: "reset password not found"})
	}
	if err := services.AuthenticateResetPassword(id, resetToken); err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "Edit", c).Info("cannot authenticate reset password: %v", err)
		return err
	}
	return c.Render(http.StatusOK, "edit_password.html.tpl", map[string]interface{}{
		"title":      "PasswordReset",
		"token":      token,
		"id":         id,
		"resetToken": resetToken,
	})
}

func (u *Passwords) Update(c echo.Context) error {
	var editPasswordForm EditPasswordForm
	err := c.Bind(editPasswordForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameters")
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "Update", err, c).Error(err)
		return err
	}

	if !CheckCSRFToken(c, editPasswordForm.Token) {
		err := errors.New("cannot verify CSRF token")
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "Update", err, c).Error(err)
		return err
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordsController", "Update", err, c).Error(err)
		return c.JSON(http.StatusNotFound, &JSONError{message: "reset password not found"})
	}

	valid, err := validators.PasswordUpdateValidation(editPasswordForm.ResetToken, editPasswordForm.Password, editPasswordForm.PasswordConfirm)
	if err != nil || !valid {
		logging.SharedInstance().MethodInfo("PasswordController", "Update", c).Infof("validation failed: %v", err)
		return c.Redirect(http.StatusFound, "/passwords/"+string(id)+"/edit")
	}

	targetUser, err := handlers.ChangeUserPassword(id, editPasswordForm.ResetToken, editPasswordForm.Password)
	if err != nil {
		logging.SharedInstance().MethodInfo("PasswordsController", "Update", c).Infof("cannot authenticate reset password: %v", err)
		return err
	}

	go password_mailer.Changed(targetUser.UserEntity.UserModel.Email)
	logging.SharedInstance().MethodInfo("PasswordsController", "Update", c).Info("success to change password")
	return c.Redirect(http.StatusFound, "/sign_in")
}

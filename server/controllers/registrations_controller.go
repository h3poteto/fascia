package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/validators"

	"html/template"
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type Registrations struct {
}

type SignUpForm struct {
	Email           string `param:"email"`
	Password        string `param:"password"`
	PasswordConfirm string `param:"password_confirm"`
	Token           string `param:"token"`
}

func (u *Registrations) SignUp(c echo.Context) error {
	url := githubOauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	token, err := GenerateCSRFToken(c)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("RegistrationsController", "SignUp", err, c).Errorf("CSRF error: %v", err)
		return err
	}

	return c.Render(http.StatusOK, "sign_up.html.tpl", map[string]interface{}{
		"title":    "SignUp",
		"oauthURL": url,
		"token":    token,
	})
}

func (u *Registrations) Registration(c echo.Context) error {
	var signUpForm SignUpForm
	err := c.Bind(signUpForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().MethodInfoWithStacktrace("RegistrationsController", "SignUp", err, c).Error(err)
		return err
	}
	logging.SharedInstance().MethodInfo("RegistrationsController", "SignUp", c).Debugf("post registration form: %+v", signUpForm)

	if !CheckCSRFToken(c, signUpForm.Token) {
		err := errors.New("cannot verify CSRF token")
		logging.SharedInstance().MethodInfoWithStacktrace("RegistrationsController", "SignUp", err, c).Error(err)
		return err
	}

	// sign up
	valid, err := validators.UserRegistrationValidation(signUpForm.Email, signUpForm.Password, signUpForm.PasswordConfirm)
	// TODO: 失敗していることは何かしらの方法で伝えたい
	if err != nil || !valid {
		logging.SharedInstance().MethodInfo("RegistrationsController", "SignUp", c).Infof("validation failed: %v", err)
		return c.Redirect(http.StatusFound, "/sign_up")
	}
	// TODO: ここCSRFのmiddlewareとかでなんとかならんかなぁ
	_, err = handlers.RegistrationUser(
		template.HTMLEscapeString(signUpForm.Email),
		template.HTMLEscapeString(signUpForm.Password),
		template.HTMLEscapeString(signUpForm.PasswordConfirm),
	)
	if err != nil {
		// TODO: 登録情報が間違っていることを通知したい
		logging.SharedInstance().MethodInfo("RegistrationsController", "SignUp", c).Infof("registration error: %v", err)
		return c.Redirect(http.StatusFound, "/sign_up")
	}

	// TODO: 成功していることも伝えたい
	logging.SharedInstance().MethodInfo("RegistrationsController", "SignUp", c).Info("registration success")
	return c.Redirect(http.StatusFound, "/sign_in")
}

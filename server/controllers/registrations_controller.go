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

// Registrations is controlelr struct for registrations
type Registrations struct {
}

// SignUpForm is struct for sign up
type SignUpForm struct {
	Email           string `json:"email" form:"email"`
	Password        string `json:"password" form:"password"`
	PasswordConfirm string `json:"password_confirm" form:"password_confirm"`
	Token           string `json:"token" form:"token"`
}

// SignUp render sign up form
func (u *Registrations) SignUp(c echo.Context) error {
	privateURL := githubPrivateConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	publicURL := githubPublicConf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	token, err := GenerateCSRFToken(c)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Errorf("CSRF error: %v", err)
		return err
	}

	return c.Render(http.StatusOK, "sign_up.html.tpl", map[string]interface{}{
		"title":      "SignUp",
		"privateURL": privateURL,
		"publicURL":  publicURL,
		"token":      token,
	})
}

// Registration creates a new user
func (u *Registrations) Registration(c echo.Context) error {
	signUpForm := new(SignUpForm)
	if err := c.Bind(signUpForm); err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Debugf("post registration form: %+v", signUpForm)

	if !CheckCSRFToken(c, signUpForm.Token) {
		err := errors.New("cannot verify CSRF token")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	// sign up
	valid, err := validators.UserRegistrationValidation(signUpForm.Email, signUpForm.Password, signUpForm.PasswordConfirm)
	// TODO: 失敗していることは何かしらの方法で伝えたい
	if err != nil || !valid {
		logging.SharedInstance().Controller(c).Infof("validation failed: %v", err)
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
		logging.SharedInstance().Controller(c).Infof("registration error: %v", err)
		return c.Redirect(http.StatusFound, "/sign_up")
	}

	// TODO: 成功していることも伝えたい
	logging.SharedInstance().Controller(c).Info("registration success")
	return c.Redirect(http.StatusFound, "/sign_in")
}

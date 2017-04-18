package controllers

import (
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/handlers"

	"html/template"
	"net/http"
	"time"

	"github.com/ipfans/echo-session"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// Webviews controller struct
type Webviews struct {
}

// SignIn is a sign in action for mobile app
func (u *Webviews) SignIn(c echo.Context) error {
	url := githubOauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	token, err := GenerateCSRFToken(c)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Errorf("CSRF error: %v", err)
		return err
	}

	// prepare cookie
	cookie := http.Cookie{
		Path:    "/",
		Name:    "fascia-ios",
		Value:   "login-session",
		Expires: time.Now().AddDate(0, 0, 1),
	}
	c.SetCookie(&cookie)

	return c.Render(http.StatusOK, "webviews/sign_in.html.tpl", map[string]interface{}{
		"title":    "SignIn",
		"oauthURL": url,
		"token":    token,
	})
}

// NewSession is a sign in action for mobile app
func (u *Webviews) NewSession(c echo.Context) error {
	s := session.Default(c)
	s.Clear()
	err := s.Save()
	if err != nil {
		err := errors.Wrap(err, "session error")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	var signInForm SignInForm
	err = c.Bind(signInForm)
	if err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	if !CheckCSRFToken(c, signInForm.Token) {
		err := errors.New("cannot verify CSRF token")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	userService, err := handlers.LoginUser(template.HTMLEscapeString(signInForm.Email), template.HTMLEscapeString(signInForm.Password))
	if err != nil {
		logging.SharedInstance().Controller(c).Infof("login error: %v", err)
		return c.Redirect(http.StatusFound, "/webviews/sign_in")
	}
	logging.SharedInstance().Controller(c).Debugf("login success: %+v", userService)
	s.Options(session.Options{
		Path:   "/",
		MaxAge: config.Element("session").(map[interface{}]interface{})["timeout"].(int),
	})
	s.Set("current_user_id", userService.UserEntity.UserModel.ID)
	err = s.Save()
	if err != nil {
		err := errors.Wrap(err, "session error")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Info("login success")
	return c.Redirect(http.StatusFound, "/webviews/callback")
}

// Callback is a empty page for mobile application handling
func (u *Webviews) Callback(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

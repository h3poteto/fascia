package controllers

import (
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/handlers"

	"html/template"
	"net/http"

	"github.com/ipfans/echo-session"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type Sessions struct {
}

type SignInForm struct {
	Email    string `param:"email"`
	Password string `param:"password"`
	Token    string `param:"token"`
}

func (u *Sessions) SignIn(c echo.Context) error {
	url := githubOauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	token, err := GenerateCSRFToken(c)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Errorf("CSRF error: %v", err)
		return err
	}

	return c.Render(http.StatusOK, "sign_in.html.tpl", map[string]interface{}{
		"title":    "SignIn",
		"oauthURL": url,
		"token":    token,
	})
}

func (u *Sessions) NewSession(c echo.Context) error {
	// 旧セッションの削除
	s := session.Default(c)
	s.Clear()
	s.Save()

	var signInForm SignInForm
	err := c.Bind(signInForm)
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
		return c.Redirect(http.StatusFound, "/sign_in")
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
	return c.Redirect(http.StatusFound, "/")
}

func (u *Sessions) SignOut(c echo.Context) error {
	s := session.Default(c)
	s.Clear()
	err := s.Save()
	if err != nil {
		err := errors.Wrap(err, "session error")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Info("logout success")
	return c.Redirect(http.StatusFound, "/sign_in")
}

func (u *Sessions) Update(c echo.Context) error {
	userService, err := LoginRequired(c)
	if err != nil {
		logging.SharedInstance().Controller(c).Infof("login error: %v", err)
		return NewJSONError(err, http.StatusUnauthorized, c)
	}
	logging.SharedInstance().Controller(c).Info("login success")

	s := session.Default(c)
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
	logging.SharedInstance().Controller(c).Info("session update success")
	return c.JSON(http.StatusOK, nil)
}

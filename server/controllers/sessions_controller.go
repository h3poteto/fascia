package controllers

import (
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/middlewares"
	"github.com/h3poteto/fascia/server/session"

	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Sessions is controller struct for sessions
type Sessions struct {
}

// SignIn renders a sign in form
func (u *Sessions) SignIn(c echo.Context) error {
	token, err := GenerateCSRFToken(c)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Errorf("CSRF error: %v", err)
		return err
	}

	return c.Render(http.StatusOK, "sign_in.html.tpl", map[string]interface{}{
		"title": "SignIn",
		"token": token,
	})
}

// SignOut delete a session and logout
func (u *Sessions) SignOut(c echo.Context) error {
	err := session.SharedInstance().Clear(c.Request(), c.Response())
	if err != nil {
		err := errors.Wrap(err, "session error")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	logging.SharedInstance().Controller(c).Info("logout success")
	return c.JSON(http.StatusOK, nil)
}

// Update a session
func (u *Sessions) Update(c echo.Context) error {
	uc, ok := c.(*middlewares.LoginContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	userService := uc.CurrentUser

	option := &sessions.Options{
		Path:     "/",
		MaxAge:   config.Element("session").(map[interface{}]interface{})["timeout"].(int),
		HttpOnly: true,
	}
	err := session.SharedInstance().Set(c.Request(), c.Response(), "current_user_id", userService.ID, option)
	if err != nil {
		err := errors.Wrap(err, "session error")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Info("session update success")
	return c.JSON(http.StatusOK, nil)
}

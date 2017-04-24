package controllers

import (
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/session"

	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// Oauth is controller struct for oauth
type Oauth struct {
}

// Github catch callback from github for oauth login
func (u *Oauth) Github(c echo.Context) error {
	// 旧セッションの削除
	err := session.SharedInstance().Clear(c.Request(), c.Response())
	if err != nil {
		err := errors.Wrap(err, "session error")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	code := c.QueryParam("code")
	logging.SharedInstance().Controller(c).Debugf("github callback param: %+v", code)
	token, err := githubOauthConf.Exchange(oauth2.NoContext, code)
	logging.SharedInstance().Controller(c).Debugf("token: %v", token)
	if err != nil {
		err := errors.Wrap(err, "oauth token error")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	userService, err := handlers.FindOrCreateUserFromGithub(token.AccessToken)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return c.Redirect(http.StatusFound, "/sign_in")
	}
	logging.SharedInstance().Controller(c).Debugf("login success: %+v", userService)

	option := &sessions.Options{Path: "/", MaxAge: config.Element("session").(map[interface{}]interface{})["timeout"].(int)}
	err = session.SharedInstance().Set(c.Request(), c.Response(), "current_user_id", userService.UserEntity.UserModel.ID, option)
	if err != nil {
		err := errors.Wrap(err, "session error")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Info("github login success")

	// iosからのセッションの場合はリダイレクト先を変更
	cookie, err := c.Cookie("fascia-ios")
	if err == nil && cookie.Value == "login-session" {
		return c.Redirect(http.StatusFound, "/webviews/callback")
	}
	return c.Redirect(http.StatusFound, "/")
}

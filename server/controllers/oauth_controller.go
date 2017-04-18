package controllers

import (
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/handlers"

	"net/http"

	"github.com/ipfans/echo-session"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type Oauth struct {
}

func (u *Oauth) Github(c echo.Context) error {
	// 旧セッションの削除
	s := session.Default(c)
	s.Clear()
	s.Save()

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
	logging.SharedInstance().Controller(c).Info("github login success")

	// iosからのセッションの場合はリダイレクト先を変更
	cookie, err := c.Cookie("fascia-ios")
	if err == nil && cookie.Value == "login-session" {
		return c.Redirect(http.StatusFound, "/webviews/callback")
	}
	return c.Redirect(http.StatusFound, "/")
}

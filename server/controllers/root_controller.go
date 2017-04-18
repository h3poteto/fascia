package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/handlers"

	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"golang.org/x/oauth2"
)

type Root struct {
}

func (u *Root) Index(c echo.Context) error {
	currentUser, err := LoginRequired(c)
	// ログインしていない場合はaboutページを見せる
	if err != nil {
		logging.SharedInstance().MethodInfo("RootController", "Index", c).Infof("login error: %v", err)
		return u.About(c)
	}
	// ログインしている場合はダッシュボードへ
	logging.SharedInstance().MethodInfo("RootController", "Index", c).Info("login success")

	projectID, _ := strconv.ParseInt(c.Param("project_id"), 10, 64)
	if projectID != 0 {
		projectService, err := handlers.FindProject(projectID)
		if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
			logging.SharedInstance().MethodInfo("RootController", "Index", c).Warnf("project not found: %v", err)
			return NewJSONError(err, http.StatusNotFound, c)
		}
	}
	return c.Render(http.StatusOK, "home.html.tpl", map[string]interface{}{
		"title": "Fascia",
	})
}

func (u *Root) About(c echo.Context) error {
	url := githubOauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	token, err := GenerateCSRFToken(c)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("RootController", "About", err, c).Error(err)
		return err
	}
	return c.Render(http.StatusOK, "about.html.tpl", map[string]interface{}{
		"title":    "Fascia",
		"oauthURL": url,
		"token":    token,
	})
}

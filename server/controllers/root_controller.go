package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/middlewares"

	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"golang.org/x/oauth2"
)

// Root is controller struct
type Root struct {
}

// Index render a top page
func (u *Root) Index(c echo.Context) error {
	currentUser, err := middlewares.CheckLogin(c)
	// ログインしていない場合はaboutページを見せる
	if err != nil {
		logging.SharedInstance().Controller(c).Infof("login error: %v", err)
		return u.About(c)
	}
	// ログインしている場合はダッシュボードへ
	logging.SharedInstance().Controller(c).Info("login success")

	projectID, _ := strconv.ParseInt(c.Param("project_id"), 10, 64)
	if projectID != 0 {
		projectService, err := handlers.FindProject(projectID)
		if err != nil || !(projectService.CheckOwner(currentUser.UserEntity.UserModel.ID)) {
			logging.SharedInstance().Controller(c).Warnf("project not found: %v", err)
			return NewJSONError(err, http.StatusNotFound, c)
		}
	}
	return c.Render(http.StatusOK, "home.html.tpl", map[string]interface{}{
		"title": "Fascia",
	})
}

// About render a about
func (u *Root) About(c echo.Context) error {
	privateURL := githubPrivateConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	publicURL := githubPublicConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	token, err := GenerateCSRFToken(c)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	return c.Render(http.StatusOK, "about.html.tpl", map[string]interface{}{
		"title":      "Fascia",
		"privateURL": privateURL,
		"publicURL":  publicURL,
		"token":      token,
	})
}

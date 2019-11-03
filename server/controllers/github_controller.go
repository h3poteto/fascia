package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/middlewares"
	"github.com/h3poteto/fascia/server/usecases/board"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Github is controller struct for github
type Github struct {
}

// Repositories returns github repositories
func (u *Github) Repositories(c echo.Context) error {
	uc, ok := c.(*middlewares.LoginContext)
	if !ok {
		err := errors.New("Can not cast context")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	currentUser := uc.CurrentUser
	if !currentUser.OauthToken.Valid {
		logging.SharedInstance().Controller(c).Info("user did not have oauth")
		return c.JSON(http.StatusOK, nil)
	}

	repositories, err := board.GetAllRepositories(currentUser.OauthToken.String)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	logging.SharedInstance().Controller(c).Info("success to get repositories")
	return c.JSON(http.StatusOK, repositories)
}

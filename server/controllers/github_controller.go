package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/middlewares"

	"context"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
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
	currentUser := uc.CurrentUserService
	if !currentUser.UserEntity.UserModel.OauthToken.Valid {
		logging.SharedInstance().Controller(c).Info("user did not have oauth")
		return c.JSON(http.StatusOK, nil)
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: currentUser.UserEntity.UserModel.OauthToken.String},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	nextPage := -1
	var repositories []*github.Repository
	for nextPage != 0 {
		if nextPage < 0 {
			nextPage = 0
		}
		repositoryOption := &github.RepositoryListOptions{
			Type:      "all",
			Sort:      "full_name",
			Direction: "asc",
			ListOptions: github.ListOptions{
				Page:    nextPage,
				PerPage: 50,
			},
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		repos, res, err := client.Repositories.List(ctx, "", repositoryOption)
		nextPage = res.NextPage
		if err != nil {
			err := errors.Wrap(err, "repository error")
			logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
			return err
		}
		repositories = append(repositories, repos...)

	}
	logging.SharedInstance().Controller(c).Info("success to get repositories")
	return c.JSON(http.StatusOK, repositories)
}

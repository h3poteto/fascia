package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/usecases/board"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Repositories defines repositories_controller methods
type Repositories struct {
}

// Hook catche events from github
func (u *Repositories) Hook(c echo.Context) error {
	eventType := c.Request().Header.Get("X-GitHub-Event")
	signature := c.Request().Header.Get("X-Hub-Signature")
	deliveryID := c.Request().Header.Get("X-GitHub-Delivery")

	if eventType == "" || signature == "" || deliveryID == "" {
		// TODO: ここエラーにしてslcak通知ほしいかも
		logging.SharedInstance().Controller(c).Errorf("could not find header information: %+v", c.Request().Header)
		return NewJSONError(errors.New("event, signature, or delivery_id is not exist"), http.StatusNotFound, c)
	}

	switch eventType {
	case "issues":
		var githubBody github.IssuesEvent
		data, _ := ioutil.ReadAll(c.Request().Body)
		json.Unmarshal(data, &githubBody)
		id := int64(*githubBody.Repo.ID)

		repo, err := board.FindRepositoryByGithubRepoID(id)
		if err != nil {
			logging.SharedInstance().ControllerWithStacktrace(err, c).Errorf("could not find repository: %v", err)
			return NewJSONError(err, http.StatusNotFound, c)
		}
		if err := repo.Authenticate(signature, data); err != nil {
			logging.SharedInstance().Controller(c).Infof("cannot authenticate to repository: %v", err)
			return NewJSONError(err, http.StatusNotFound, c)
		}
		err = board.ApplyIssueChangesToRepository(repo, githubBody)
		if err != nil {
			logging.SharedInstance().ControllerWithStacktrace(err, c).Errorf("could not apply issue changes: %v", err)
			return err
		}
		logging.SharedInstance().Controller(c).Info("success apply issues event from webhook")

	case "pull_request":
		var githubBody github.PullRequestEvent
		data, _ := ioutil.ReadAll(c.Request().Body)
		json.Unmarshal(data, &githubBody)
		id := int64(*githubBody.Repo.ID)

		repo, err := board.FindRepositoryByGithubRepoID(id)
		if err != nil {
			logging.SharedInstance().ControllerWithStacktrace(err, c).Errorf("could not find repository: %v", err)
			return NewJSONError(err, http.StatusNotFound, c)
		}
		if err := repo.Authenticate(signature, data); err != nil {
			logging.SharedInstance().Controller(c).Infof("cannot authenticate to repository: %v", err)
			return NewJSONError(err, http.StatusNotFound, c)
		}

		err = board.ApplyPullRequestChangesToRepository(repo, githubBody)
		if err != nil {
			logging.SharedInstance().ControllerWithStacktrace(err, c).Errorf("could not apply pull request changes: %v", err)
			return err
		}
		logging.SharedInstance().Controller(c).Info("success apply pull request event from webhook")
	}

	return c.JSON(http.StatusOK, nil)
}

package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/handlers"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/labstack/echo"
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
		logging.SharedInstance().MethodInfo("Repositories", "Hook", c).Infof("could not find header information: %+v", c.Request().Header)
		return c.JSON(http.StatusNotFound, &JSONError{message: "event, signature, or delivery_id is not exist"})
	}

	switch eventType {
	case "issues":
		var githubBody github.IssuesEvent
		data, _ := ioutil.ReadAll(c.Request().Body)
		json.Unmarshal(data, &githubBody)
		id := int64(*githubBody.Repo.ID)

		repo, err := handlers.FindRepositoryByGithubRepoID(id)
		if err != nil {
			logging.SharedInstance().MethodInfoWithStacktrace("Repositories", "Hook", err, c).Errorf("could not find repository: %v", err)
			return c.JSON(http.StatusNotFound, &JSONError{message: "repository not found"})
		}
		if err := repo.Authenticate(signature, data); err != nil {
			logging.SharedInstance().MethodInfo("Repositories", "Hook", c).Infof("cannot authenticate to repository: %v", err)
			return c.JSON(http.StatusNotFound, &JSONError{message: "repository not found"})
		}
		err = handlers.ApplyIssueChangesToRepository(repo, githubBody)
		if err != nil {
			logging.SharedInstance().MethodInfoWithStacktrace("Repositories", "Hook", err, c).Error("could not apply issue changes: %v", err)
			return err
		}
		logging.SharedInstance().MethodInfo("Repositories", "Hook", c).Info("success apply issues event from webhook")

	case "pull_request":
		var githubBody github.PullRequestEvent
		data, _ := ioutil.ReadAll(c.Request().Body)
		json.Unmarshal(data, &githubBody)
		id := int64(*githubBody.Repo.ID)

		repo, err := handlers.FindRepositoryByGithubRepoID(id)
		if err != nil {
			logging.SharedInstance().MethodInfoWithStacktrace("Repositories", "Hook", err, c).Errorf("could not find repository: %v", err)
			return c.JSON(http.StatusNotFound, &JSONError{message: "repository not found"})
		}
		if err := repo.Authenticate(signature, data); err != nil {
			logging.SharedInstance().MethodInfo("Repositories", "Hook", c).Infof("cannot authenticate to repository: %v", err)
			return c.JSON(http.StatusNotFound, &JSONError{message: "repository not found"})
		}

		err = handlers.ApplyPullRequestChangesToRepository(repo, githubBody)
		if err != nil {
			logging.SharedInstance().MethodInfoWithStacktrace("Repositories", "Hook", err, c).Errorf("could not apply pull request changes: %v", err)
			return err
		}
		logging.SharedInstance().MethodInfo("Repositories", "Hook", c).Info("success apply pull request event from webhook")
	}

	return c.JSON(http.StatusOK, nil)
}

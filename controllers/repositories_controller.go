package controllers

import (
	"../models/project"
	"../models/repository"
	"../modules/logging"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/zenazn/goji/web"
)

// Repositories defines repositories_controller methods
type Repositories struct {
}

// Hook catche events from github
func (u *Repositories) Hook(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	eventType := r.Header.Get("X-GitHub-Event")
	signature := r.Header.Get("X-Hub-Signature")
	deliveryID := r.Header.Get("X-GitHub-Delivery")

	if eventType == "" || signature == "" || deliveryID == "" {
		logging.SharedInstance().MethodInfo("Repositories", "Hook").Infof("could not find header information: %v", r.Header)
		return
	}

	switch eventType {
	case "issues":
		var githubBody github.IssuesEvent
		data, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(data, &githubBody)
		id := int64(*githubBody.Repo.ID)

		repo := repository.FindRepositoryByRepositoryID(id)
		if repo == nil {
			logging.SharedInstance().MethodInfo("Repositories", "Hook").Error("could not find repository")
			return
		}
		if !repo.Authenticate(signature, data) {
			logging.SharedInstance().MethodInfo("Repositories", "Hook").Info("cannot authenticate to repository")
			return
		}
		project.IssuesEvent(repo.ID, githubBody)

	case "pull_request":
		var githubBody github.PullRequestEvent
		data, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(data, &githubBody)
		id := int64(*githubBody.Repo.ID)

		repo := repository.FindRepositoryByRepositoryID(id)
		if repo == nil {
			logging.SharedInstance().MethodInfo("Repositories", "Hook").Error("could not find repository")
			return
		}
		if !repo.Authenticate(signature, data) {
			logging.SharedInstance().MethodInfo("Repositories", "Hook").Info("cannot authenticate to repository")
			return
		}
		project.PullRequestEvent(repo.ID, githubBody)
	}
}

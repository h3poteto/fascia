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
		http.Error(w, "event, signature, or delivery_id is not exist", 404)
		return
	}

	switch eventType {
	case "issues":
		var githubBody github.IssuesEvent
		data, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(data, &githubBody)
		id := int64(*githubBody.Repo.ID)

		repo, err := repository.FindRepositoryByRepositoryID(id)
		if err != nil {
			logging.SharedInstance().MethodInfo("Repositories", "Hook", true).Errorf("could not find repository: %v", err)
			http.Error(w, "repository not found", 404)
			return
		}
		if err := repo.Authenticate(signature, data); err != nil {
			logging.SharedInstance().MethodInfo("Repositories", "Hook").Infof("cannot authenticate to repository: %v", err)
			http.Error(w, "repository authenticate failed", 404)
			return
		}
		err = project.IssuesEvent(repo.ID, githubBody)
		if err != nil {
			logging.SharedInstance().MethodInfo("Repositories", "Hook", true).Errorf("cannot handle issue event: %v", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}

	case "pull_request":
		var githubBody github.PullRequestEvent
		data, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(data, &githubBody)
		id := int64(*githubBody.Repo.ID)

		repo, err := repository.FindRepositoryByRepositoryID(id)
		if err != nil {
			logging.SharedInstance().MethodInfo("Repositories", "Hook", true).Errorf("could not find repository: %v", err)
			http.Error(w, "repository not found", 404)
			return
		}
		if err := repo.Authenticate(signature, data); err != nil {
			logging.SharedInstance().MethodInfo("Repositories", "Hook").Infof("cannot authenticate to repository: %v", err)
			http.Error(w, "repository authenticate failed", 404)
			return
		}
		err = project.PullRequestEvent(repo.ID, githubBody)
		if err != nil {
			logging.SharedInstance().MethodInfo("Repositories", "Hook", true).Errorf("cannot handle pull request event: %v", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
	}

	return
}

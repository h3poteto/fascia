package handlers

import (
	"math/rand"
	"time"

	"github.com/google/go-github/github"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/services"
)

// ApplyIssueChangesToRepository apply updating information of issue to each task
func ApplyIssueChangesToRepository(repository *services.Repository, githubBody github.IssuesEvent) error {
	projectServices, err := services.FindProjectByRepositoryID(repository.RepositoryEntity.ID)
	if err != nil {
		return err
	}

	go func(projectServices []*services.Project, githubBody github.IssuesEvent) {
		waitWebhookReadtime()
		for _, p := range projectServices {
			err = p.ApplyIssueChanges(githubBody)
			if err != nil {
				logging.SharedInstance().MethodInfoWithStacktrace("Webhook", "ApplyIssueChangesToRepository", err).Error(err)
				return
			}
		}
	}(projectServices, githubBody)
	return nil
}

// ApplyPullRequestChangesToRepository apply updating information of pull request to each task
func ApplyPullRequestChangesToRepository(repository *services.Repository, githubBody github.PullRequestEvent) error {

	projectServices, err := services.FindProjectByRepositoryID(repository.RepositoryEntity.ID)
	if err != nil {
		return err
	}

	go func(projectServices []*services.Project, githubBody github.PullRequestEvent) {
		waitWebhookReadtime()
		for _, p := range projectServices {
			err = p.ApplyPullRequestChanges(githubBody)
			if err != nil {
				logging.SharedInstance().MethodInfoWithStacktrace("Webhook", "ApplyPullRequestChangesToRepository", err).Error(err)
				return
			}
		}
	}(projectServices, githubBody)
	return nil
}

func waitWebhookReadtime() {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration((rand.Intn(20) + 5)) * time.Second)
}

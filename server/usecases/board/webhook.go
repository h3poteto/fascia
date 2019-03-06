package board

import (
	"math/rand"
	"time"

	"github.com/google/go-github/github"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/domains/project"
	"github.com/h3poteto/fascia/server/domains/repo"
)

// ApplyIssueChangesToRepository apply updating information of issue to each task
func ApplyIssueChangesToRepository(repo *repo.Repo, githubBody github.IssuesEvent) error {
	projects, err := findProjectByRepoID(repo.ID)
	if err != nil {
		return err
	}

	go func(projects []*project.Project, githubBody github.IssuesEvent) {
		waitWebhookReadtime()
		for _, p := range projects {
			err = applyIssueChanges(p, githubBody)
			if err != nil {
				logging.SharedInstance().MethodInfoWithStacktrace("Webhook", "ApplyIssueChangesToRepository", err).Error(err)
				return
			}
		}
	}(projects, githubBody)
	return nil
}

// ApplyPullRequestChangesToRepository apply updating information of pull request to each task
func ApplyPullRequestChangesToRepository(repo *repo.Repo, githubBody github.PullRequestEvent) error {

	projects, err := findProjectByRepoID(repo.ID)
	if err != nil {
		return err
	}

	go func(projects []*project.Project, githubBody github.PullRequestEvent) {
		waitWebhookReadtime()
		for _, p := range projects {
			err = applyPullRequestChanges(p, githubBody)
			if err != nil {
				logging.SharedInstance().MethodInfoWithStacktrace("Webhook", "ApplyPullRequestChangesToRepository", err).Error(err)
				return
			}
		}
	}(projects, githubBody)
	return nil
}

func waitWebhookReadtime() {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration((rand.Intn(20) + 5)) * time.Second)
}

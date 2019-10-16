package services

import (
	"math/rand"
	"time"

	"github.com/google/go-github/github"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/domains/list"
	"github.com/h3poteto/fascia/server/domains/project"
	"github.com/h3poteto/fascia/server/domains/repo"
	"github.com/h3poteto/fascia/server/domains/task"
)

func ApplyIssueChangesToRepository(projects []*project.Project, githubBody github.IssuesEvent, projectInfra project.Repository, listInfra list.Repository, taskInfra task.Repository, repoInfra repo.Repository) {
	waitWebhookReadtime()
	for _, p := range projects {
		err := applyIssueChanges(p, githubBody, projectInfra, listInfra, taskInfra, repoInfra)
		if err != nil {
			logging.SharedInstance().MethodInfoWithStacktrace("Webhook", "ApplyIssueChangesToRepository", err).Error(err)
			return
		}
	}

}

func waitWebhookReadtime() {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration((rand.Intn(20) + 5)) * time.Second)
}

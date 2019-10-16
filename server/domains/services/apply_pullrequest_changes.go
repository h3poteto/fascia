package services

import (
	"github.com/google/go-github/github"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/domains/list"
	"github.com/h3poteto/fascia/server/domains/project"
	"github.com/h3poteto/fascia/server/domains/repo"
	"github.com/h3poteto/fascia/server/domains/task"
)

func ApplyPullRequestChangesToRepository(projects []*project.Project, githubBody github.PullRequestEvent, projectInfra project.Repository, listInfra list.Repository, taskInfra task.Repository, repoInfra repo.Repository) {
	waitWebhookReadtime()
	for _, p := range projects {
		err := applyPullRequestChanges(p, githubBody, projectInfra, listInfra, taskInfra, repoInfra)
		if err != nil {
			logging.SharedInstance().MethodInfoWithStacktrace("Webhook", "ApplyPullRequestChangesToRepository", err).Error(err)
			return
		}
	}
}

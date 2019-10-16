package board

import (
	"github.com/google/go-github/github"
	"github.com/h3poteto/fascia/server/domains/repo"
	"github.com/h3poteto/fascia/server/domains/services"
)

// ApplyIssueChangesToRepository apply updating information of issue to each task
func ApplyIssueChangesToRepository(repo *repo.Repo, githubBody github.IssuesEvent) error {
	projectInfra := InjectProjectRepository()
	projects, err := projectInfra.FindByRepositoryID(repo.ID)
	if err != nil {
		return err
	}

	go services.ApplyIssueChangesToRepository(projects, githubBody, InjectProjectRepository(), InjectListRepository(), InjectTaskRepository(), InjectRepoRepository())
	return nil
}

// ApplyPullRequestChangesToRepository apply updating information of pull request to each task
func ApplyPullRequestChangesToRepository(repo *repo.Repo, githubBody github.PullRequestEvent) error {
	projectInfra := InjectProjectRepository()
	projects, err := projectInfra.FindByRepositoryID(repo.ID)
	if err != nil {
		return err
	}

	go services.ApplyPullRequestChangesToRepository(projects, githubBody, InjectProjectRepository(), InjectListRepository(), InjectTaskRepository(), InjectRepoRepository())
	return nil
}

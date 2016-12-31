package handlers

import (
	"strings"

	"github.com/google/go-github/github"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/services"
	"github.com/pkg/errors"
)

// FindRepositoryByGithubRepoID search repository according to github repository id
func FindRepositoryByGithubRepoID(id int64) (*services.Repository, error) {
	return services.FindRepositoryByGithubRepoID(id)
}

func ApplyIssueChangesToRepository(repository *services.Repository, githubBody github.IssuesEvent) error {
	projectServices, err := services.FindProjectByRepositoryID(repository.RepositoryEntity.RepositoryModel.ID)
	if err != nil {
		return err
	}

	for _, p := range projectServices {
		err = p.ApplyIssueChanges(githubBody)
		if err != nil {
			if !includeDuplicateError(err) {
				return err
			}
			logging.SharedInstance().MethodInfo("Handlers", "ApplyIssueChangesToRepository").Warn(err)
		}
	}
	return nil
}

func ApplyPullRequestChangesToRepository(repository *services.Repository, githubBody github.PullRequestEvent) error {

	projectServices, err := services.FindProjectByRepositoryID(repository.RepositoryEntity.RepositoryModel.ID)
	if err != nil {
		return err
	}
	for _, p := range projectServices {
		err = p.ApplyPullRequestChanges(githubBody)
		if err != nil {
			if !includeDuplicateError(err) {
				return err
			}
			logging.SharedInstance().MethodInfo("Handlers", "ApplyPullRequestChangesToRepository").Warn(err)
		}
	}
	return nil
}

func includeDuplicateError(err error) bool {
	if strings.Index(errors.Cause(err).Error(), "Error 1062") == 0 {
		return true
	}
	return false
}

package handlers

import (
	"database/sql"

	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/services"
)

func CreateProject(userID int64, title string, description string, repositoryID int, oauthToken sql.NullString) (*services.Project, error) {
	projectService := services.NewProject(nil)
	_, err := projectService.Create(userID, title, description, repositoryID, oauthToken)
	if err != nil {
		return nil, err
	}

	go func(projectService *services.Project) {
		// Create initial list before get issues from github
		err := projectService.FetchCreatedInitialList()
		if err != nil {
			return
		}
		// Sync issues from github
		_, err = projectService.FetchGithub()
		if err != nil {
			logging.SharedInstance().MethodInfoWithStacktrace("Project", "Create", err).Error(err)
			return
		}

		// Create Webhook in github
		err = projectService.CreateWebhook()
		if err != nil {
			logging.SharedInstance().MethodInfo("Project", "Create").Infof("failed to create webhook: %v", err)
			return
		}
		logging.SharedInstance().MethodInfo("Project", "Create").Info("success to create webhook")
	}(projectService)

	return projectService, nil
}

func FindProject(projectID int64) (*services.Project, error) {
	return services.FindProject(projectID)
}

func FindProjectByRepositoryID(repositoryID int64) (*services.Project, error) {
	return services.FindProjectByRepositoryID(repositoryID)
}

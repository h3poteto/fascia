package handlers

import (
	"database/sql"

	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/services"
)

// CreateProject create a new project, and fetch github
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

// FindProject search a project according to project id
func FindProject(projectID int64) (*services.Project, error) {
	return services.FindProject(projectID)
}

// FindProjectByRepositoryID search a project according to repository id
func FindProjectByRepositoryID(repositoryID int64) ([]*services.Project, error) {
	return services.FindProjectByRepositoryID(repositoryID)
}

// DestroyProject delete project and delete webhook
func DestroyProject(projectID int64) error {
	projectService, err := FindProject(projectID)
	if err != nil {
		return err
	}

	// 存在しない場合は空振るので問題ない
	err = projectService.DeleteWebhook()
	if err != nil {
		return err
	}
	// repositoryと関連づいていない場合は単に削除するだけで良い
	err = projectService.DeleteLists()
	if err != nil {
		return err
	}
	err = projectService.Delete()
	if err != nil {
		return err
	}
	return nil
}

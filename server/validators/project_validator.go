package validators

import (
	"github.com/asaskevich/govalidator"
)

type projectCreate struct {
	Title        string `json:"title" valid:"required~title is required,stringlength(1|255)~title must be between 1 to 255"`
	Description  string `json:"description" valid:"stringlength(0|255)~description must be between 0 to 255,optional"`
	RepositoryID int    `json:"repository_id" valid:"-"`
}

type projectUpdate struct {
	Title       string `json:"title" valid:"requred~title is required,stringlength(1|255)~title must be between 1 to 255"`
	Description string `json:"description" valid:"stringlength(0|255)~description must be between 0 to 255,optional"`
}

// ProjectCreateValidation check form variable when create projects
func ProjectCreateValidation(title string, description string, repositoryID int) (bool, error) {
	form := &projectCreate{
		Title:        title,
		Description:  description,
		RepositoryID: repositoryID,
	}
	return govalidator.ValidateStruct(form)
}

// ProjectUpdateValidation check form variable when update projects
func ProjectUpdateValidation(title string, description string) (bool, error) {
	form := &projectUpdate{
		Title:       title,
		Description: description,
	}
	return govalidator.ValidateStruct(form)
}

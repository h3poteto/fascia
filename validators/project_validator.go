package validators

import (
	"github.com/asaskevich/govalidator"
)

type projectCreate struct {
	Title           string `valid:"stringlength(1|255)"`
	Description     string `valid:"stringlength(0|255),optional"`
	RepositoryID    int64  `valid:"int,optional"`
	RepositoryOwner string `valid:"-"`
	RepositoryName  string `valid:"-"`
}

type projectUpdate struct {
	Title       string `valid:"stringlength(1|255)"`
	Description string `valid:"stringlength(0|255),optional"`
}

func ProjectCreateValidation(title string, description string, repositoryID int64, repositoryOwner string, repositoryName string) (bool, error) {
	form := &projectCreate{
		Title:           title,
		Description:     description,
		RepositoryID:    repositoryID,
		RepositoryOwner: repositoryOwner,
		RepositoryName:  repositoryName,
	}
	return govalidator.ValidateStruct(form)
}

func ProjectUpdateValidation(title string, description string) (bool, error) {
	form := &projectUpdate{
		Title:       title,
		Description: description,
	}
	return govalidator.ValidateStruct(form)
}

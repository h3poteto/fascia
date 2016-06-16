package validators

type projectCreate struct {
	Title           string `valid:"min=1,max=255"`
	Description     string `valid:"min=0,max=255"`
	RepositoryID    int64  `valid:""`
	RepositoryOwner string `valid:""`
	RepositoryName  string `valid:""`
}

type projectUpdate struct {
	Title       string `valid:"min=1,max255"`
	Description string `valid:"min=0,max=255"`
}

func ProjectCreateValidation(title string, description string, repositoryID int64, repositoryOwner string, repositoryName string) error {
	form := &projectCreate{
		Title:           title,
		Description:     description,
		RepositoryID:    repositoryID,
		RepositoryOwner: repositoryOwner,
		RepositoryName:  repositoryName,
	}
	return validate.Struct(form)
}

func ProjectUpdateValidation(title string, description string) error {
	form := &projectUpdate{
		Title:       title,
		Description: description,
	}
	return validate.Struct(form)
}

package validators

import (
	"github.com/asaskevich/govalidator"
)

type taskCreate struct {
	Title       string `valid:"stringlength(1|255)"`
	Description string `valid:"stringlength(0|21845),optional"`
}

type taskMove struct {
	ToListID     int64 `valid:"required"`
	PrevToTaskID int64 `valid:"-"`
}

// TaskCreateValidation check form variable when create tasks
func TaskCreateValidation(title string, description string) (bool, error) {
	form := &taskCreate{
		Title:       title,
		Description: description,
	}
	return govalidator.ValidateStruct(form)
}

// TaskMoveValidation check form variable when move a task
func TaskMoveValidation(toListID int64, prevToTaskID int64) (bool, error) {
	form := &taskMove{
		ToListID:     toListID,
		PrevToTaskID: prevToTaskID,
	}
	return govalidator.ValidateStruct(form)
}

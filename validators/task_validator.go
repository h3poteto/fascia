package validators

import (
	"github.com/asaskevich/govalidator"
)

type taskCreate struct {
	Title       string `valid:"stringlength(1|255)"`
	Description string `valid:"stringlength(0|255),optional"`
}

type taskMove struct {
	ToListID     int64 `valid:"required"`
	PrevToTaskID int64 `valid:"-"`
}

func TaskCreateValidation(title string, description string) (bool, error) {
	form := &taskCreate{
		Title:       title,
		Description: description,
	}
	return govalidator.ValidateStruct(form)
}

func TaskMoveValidation(toListID int64, prevToTaskID int64) (bool, error) {
	form := &taskMove{
		ToListID:     toListID,
		PrevToTaskID: prevToTaskID,
	}
	return govalidator.ValidateStruct(form)
}

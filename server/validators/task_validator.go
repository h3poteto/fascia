package validators

import (
	"github.com/asaskevich/govalidator"
)

type taskCreate struct {
	Title       string `json:"title" valid:"required~title is required,stringlength(1|255)~title must be between 1 to 255"`
	Description string `json:"description" valid:"stringlength(0|21845)~description must be between 0 to 21845,optional"`
}

type taskUpdate struct {
	Title       string `json:"title" valid:"required~title is required,stringlength(1|255)~title must be between 1 to 255"`
	Description string `json:"description" valid:"stringlength(0|21845)~description must be between 0 to 21845,optional"`
}

type taskMove struct {
	ToListID     int64 `json:"to_list_id" valid:"required~to_list_id is required"`
	PrevToTaskID int64 `json:"prev_to_Task_id" valid:"-"`
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

// TaskUpdateValidation check form variable when update a task
func TaskUpdateValidation(title string, description string) (bool, error) {
	form := &taskUpdate{
		Title:       title,
		Description: description,
	}
	return govalidator.ValidateStruct(form)
}

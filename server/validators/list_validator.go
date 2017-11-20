package validators

import (
	"github.com/asaskevich/govalidator"
)

type listCreate struct {
	Title string `json:"title" valid:"required~title is required,stringlength(1|255)~title must be between 1 to 255"`
	Color string `json:"color" valid:"required~color is required,hexadecimal~color must be hexadecimal,stringlength(6|6)~color must be 6 characters"`
}

type listUpdate struct {
	Title    string `valid:"required~title is required,stringlength(1|255)~title must be between 1 to 255"`
	Color    string `valid:"required~color is required,hexadecimal~color must be hexadecimal,stringlength(6|6)~color must be 6 characters"`
	OptionID int64  `valid:"-"`
}

// ListCreateValidation check form variable when create lists
func ListCreateValidation(title string, color string) (bool, error) {
	form := &listCreate{
		Title: title,
		Color: color,
	}
	return govalidator.ValidateStruct(form)
}

// ListUpdateValidation check form variable when update lists
func ListUpdateValidation(title string, color string, optionID int64) (bool, error) {
	form := &listUpdate{
		Title:    title,
		Color:    color,
		OptionID: optionID,
	}
	return govalidator.ValidateStruct(form)
}

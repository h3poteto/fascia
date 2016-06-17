package validators

import (
	"github.com/asaskevich/govalidator"
)

type listCreate struct {
	Title string `valid:"stringlength(1|255)"`
	Color string `valid:"hexadecimal,stringlength=6"`
}

type listUpdate struct {
	Title  string `valid:"stringlength(1|255)"`
	Color  string `valid:"hexadecimal,stringlength(6|6)"`
	Action string `valid:"-"`
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
func ListUpdateValidation(title string, color string, action string) (bool, error) {
	form := &listUpdate{
		Title:  title,
		Color:  color,
		Action: action,
	}
	return govalidator.ValidateStruct(form)
}

package validators

type listCreate struct {
	Title string `valid:"min=1,max=255"`
	Color string `valid:"hexadecimal,len=6,required"`
}

type listUpdate struct {
	Title  string `valid:"min=1,max=255"`
	Color  string `valid:"hexadecimal,len=6,required"`
	Action string `valid:""`
}

func ListCreateValidation(title string, color string) error {
	form := &listCreate{
		Title: title,
		Color: color,
	}
	return validate.Struct(form)
}

func ListUpdateValidation(title string, color string, action string) error {
	form := &listUpdate{
		Title:  title,
		Color:  color,
		Action: action,
	}
	return validate.Struct(form)
}

package validators

import (
	"github.com/asaskevich/govalidator"
)

type inquiryCreate struct {
	Email   string `json:"email" valid:"required~email is required,stringlength(1|255)~email must be between 1 to 255"`
	Name    string `json:"name" valid:"required~name is required,stringlength(1|255)~name must be between 1 to 255"`
	Message string `json:"message" valid:"-"`
}

// InquiryCreateValidation define validation of inquiry when create a new object.
func InquiryCreateValidation(email, name, message string) (bool, error) {
	form := &inquiryCreate{
		email,
		name,
		message,
	}
	return govalidator.ValidateStruct(form)
}

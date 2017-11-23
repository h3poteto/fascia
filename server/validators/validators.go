package validators

import (
	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

// ErrorsByField call govalidator method for validation error
func ErrorsByField(err error) map[string]string {
	return govalidator.ErrorsByField(err)
}

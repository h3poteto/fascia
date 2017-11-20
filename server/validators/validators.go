package validators

import (
	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func ErrorsByField(err error) map[string]string {
	return govalidator.ErrorsByField(err)
}

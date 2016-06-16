package validators

import (
	"gopkg.in/go-playground/validator.v8"
)

var validate *validator.Validate

func init() {
	config := &validator.Config{TagName: "valid"}
	validate = validator.New(config)
}

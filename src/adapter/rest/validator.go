package rest

import "github.com/go-playground/validator/v10"

type errorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func validateInput(input interface{}) []*errorResponse {
	var errors []*errorResponse
	err := validate.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element errorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

var validate = validator.New()

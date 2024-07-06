package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}


func ValidateUser(data interface{}) []ValidationError {
    var errors []ValidationError
    validate := validator.New()

    err := validate.Struct(data)
    if err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            var element ValidationError
            element.Field = err.Field()
            switch err.Tag() {
            case "required":
                element.Message = fmt.Sprintf("%s is required", err.Field())
            case "email":
                element.Message = fmt.Sprintf("%s must be a valid email", err.Field())
            case "unique":
                element.Message = fmt.Sprintf("%s must be unique", err.Field())
            default:
                element.Message = fmt.Sprintf("%s is not valid", err.Field())
            }
            errors = append(errors, element)
        }
    }
    return errors
}

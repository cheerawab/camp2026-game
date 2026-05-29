package httpx

import (
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validateOnce sync.Once
	validate     *validator.Validate
)

func ValidateStruct(value any) error {
	v := validatorInstance()
	if err := v.Struct(value); err != nil {
		var validationErrors validator.ValidationErrors
		if ok := asValidationErrors(err, &validationErrors); ok {
			return UnprocessableEntity("invalid request body", validationErrorDetails(validationErrors)...)
		}

		return UnprocessableEntity("invalid request body")
	}

	return nil
}

func validatorInstance() *validator.Validate {
	validateOnce.Do(func() {
		validate = validator.New(validator.WithRequiredStructEnabled())
		validate.RegisterTagNameFunc(jsonTagName)
	})
	return validate
}

func jsonTagName(field reflect.StructField) string {
	name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	if name == "" {
		return field.Name
	}
	return name
}

func asValidationErrors(err error, target *validator.ValidationErrors) bool {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return false
	}
	*target = validationErrors
	return true
}

func validationErrorDetails(validationErrors validator.ValidationErrors) []ErrorDetail {
	details := make([]ErrorDetail, 0, len(validationErrors))
	for _, fieldError := range validationErrors {
		details = append(details, ErrorDetail{
			Location: validationErrorLocation(fieldError),
			Message:  validationErrorMessage(fieldError),
		})
	}
	return details
}

func validationErrorLocation(fieldError validator.FieldError) string {
	namespace := fieldError.Namespace()
	if idx := strings.Index(namespace, "."); idx >= 0 {
		namespace = namespace[idx+1:]
	}
	if namespace == "" {
		namespace = fieldError.Field()
	}
	return "body." + namespace
}

func validationErrorMessage(fieldError validator.FieldError) string {
	field := fieldError.Field()

	switch fieldError.Tag() {
	case "required":
		return field + " is required"
	case "min":
		return field + " must be at least " + fieldError.Param()
	case "max":
		return field + " must be at most " + fieldError.Param()
	case "oneof":
		return field + " must be one of: " + fieldError.Param()
	default:
		return field + " is invalid"
	}
}

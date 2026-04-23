package validation

import (
	"encoding/base64"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type UpdateProfileImageRequest struct {
	Image string `json:"image" validate:"required,base64image"`
}

func init() {
	_ = validate.RegisterValidation("base64image", validateBase64Image)
}

func ValidateProfile(s any) error {
	return validate.Struct(s)
}

func FormatValidationError(err error) map[string]string {
	errors := map[string]string{}

	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors
	}

	for _, e := range ve {
		field := toSnakeCase(e.Field())

		switch e.Tag() {
		case "required":
			errors[field] = "required"
		case "base64image":
			errors[field] = "must be a valid base64 image"
		default:
			errors[field] = "invalid"
		}
	}

	return errors
}

func validateBase64Image(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return false
	}

	if strings.Contains(value, ";base64,") {
		parts := strings.SplitN(value, ";base64,", 2)
		if len(parts) != 2 {
			return false
		}

		if !strings.HasPrefix(parts[0], "data:image/") {
			return false
		}

		value = parts[1]
	}

	_, err := base64.StdEncoding.DecodeString(value)
	return err == nil
}

func toSnakeCase(value string) string {
	if value == "" {
		return value
	}

	var result []rune
	for i, r := range value {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		if r >= 'A' && r <= 'Z' {
			r += 'a' - 'A'
		}
		result = append(result, r)
	}

	return string(result)
}

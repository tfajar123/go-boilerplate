package authValidation

import (
	"fmt"
	"go-boilerplate/ent/user"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type LoginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=72"`
}

type RegisterRequest struct {
	Name     string    `validate:"required,min=3,max=100"`
	Email    string    `validate:"required,email"`
	Password string    `validate:"required,min=8,max=72"`
	Role     user.Role `validate:"required,oneof=admin user"`
}

func ValidateAuth(s any) error {
	return validate.Struct(s)
}

func FormatValidationError(err error) map[string]string {
	errors := map[string]string{}

	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors
	}

	for _, e := range ve {
		field := strings.ToLower(e.Field())

		switch e.Tag() {

		case "required":
			errors[field] = "wajib diisi"

		case "email":
			errors[field] = "format email tidak valid"

		case "min":
			errors[field] = fmt.Sprintf("minimal %s karakter", e.Param())

		case "max":
			errors[field] = fmt.Sprintf("maksimal %s karakter", e.Param())

		case "oneof":
			errors[field] = "nilai tidak valid"

		default:
			errors[field] = "tidak valid"
		}
	}

	return errors
}

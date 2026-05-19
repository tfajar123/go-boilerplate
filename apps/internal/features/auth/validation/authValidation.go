package authValidation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type LoginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=72"`
}

type RegisterRequest struct {
	Name     string `validate:"required,min=3,max=100"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=72"`
}

type ForgotPasswordRequest struct {
	Email string `validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token           string `validate:"required"`
	NewPassword     string `validate:"required,min=8,max=72"`
	ConfirmPassword string `validate:"required,eqfield=NewPassword"`
}

type VerifyRegisterOTPRequest struct {
	Email string `validate:"required,email"`
	OTP   string `validate:"required,len=6,numeric"`
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
		field := toSnakeCase(e.Field())

		switch e.Tag() {

		case "required":
			errors[field] = "required"

		case "email":
			errors[field] = "invalid email"

		case "min":
			errors[field] = fmt.Sprintf("minimal %s characters", e.Param())

		case "max":
			errors[field] = fmt.Sprintf("maximal %s characters", e.Param())

		case "oneof":
			errors[field] = "invalid role"

		case "base64":
			errors[field] = "image must be base64"

		case "eqfield":
			errors[field] = "password confirmation does not match"

		default:
			errors[field] = "invalid"
		}
	}

	return errors
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

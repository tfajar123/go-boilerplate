package authDTO

import (
	"time"
)

// UserResponse DTO for returning user data without sensitive fields
type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginResponse DTO for successful login
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

// RefreshTokenResponse DTO for refresh token response
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ForgotPasswordResponse DTO for forgot password response
type ForgotPasswordResponse struct {
	EmailSent bool      `json:"email_sent"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ResetPasswordResponse DTO for reset password response
type ResetPasswordResponse struct {
	ResetAt time.Time `json:"reset_at"`
}

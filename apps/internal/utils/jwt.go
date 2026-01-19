package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	accessSecret  = []byte(os.Getenv("JWT_ACCESS_SECRET"))
	refreshSecret = []byte(os.Getenv("JWT_REFRESH_SECRET"))
)

// ========================
// Session
// ========================

func NewSessionID() string {
	return uuid.NewString()
}

// ========================
// Access Token
// ========================

func GenerateAccessToken(
	userID string,
	email string,
	sessionID string,
) (string, error) {

	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"sid":   sessionID,
		"type":  "access",
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}

func ParseAccessToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return accessSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("access token tidak valid")
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["type"] != "access" {
		return nil, errors.New("bukan access token")
	}

	return claims, nil
}

// ========================
// Refresh Token
// ========================

func GenerateRefreshToken(
	userID string,
	email string,
	sessionID string,
) (string, error) {

	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"sid":   sessionID,
		"type":  "refresh",
		"exp":   time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

func ParseRefreshToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return refreshSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("refresh token tidak valid")
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["type"] != "refresh" {
		return nil, errors.New("bukan refresh token")
	}

	return claims, nil
}

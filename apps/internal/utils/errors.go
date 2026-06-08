package utils

import (
	"errors"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Database Error Types
var (
	ErrDatabaseConnection  = errors.New("failed to connect to database")
	ErrDatabaseQuery       = errors.New("database query error")
	ErrNotFound            = errors.New("record not found")
	ErrDuplicateEntry      = errors.New("duplicate entry")
	ErrConstraintViolation = errors.New("constraint violation")
	ErrInvalidInput        = errors.New("invalid input data")
	ErrTimeout             = errors.New("database operation timeout")
)

// ErrorCategory represents error category for proper handling
type ErrorCategory int

const (
	CategoryUnknown ErrorCategory = iota
	CategoryNotFound
	CategoryValidation
	CategoryDatabase
	CategoryPermission
	CategoryAuthentication
	CategoryTimeout
)

// AppError is custom application error
type AppError struct {
	Err        error
	Message    string
	StatusCode int
	Category   ErrorCategory
	Details    any
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates new application error
func NewAppError(err error, message string, statusCode int, category ErrorCategory) *AppError {
	return &AppError{
		Err:        err,
		Message:    message,
		StatusCode: statusCode,
		Category:   category,
	}
}

// WrapDatabaseError wraps raw database error to proper AppError
func WrapDatabaseError(err error) *AppError {
	if err == nil {
		return nil
	}

	// MongoDB: document not found
	if errors.Is(err, mongo.ErrNoDocuments) {
		return NewAppError(err, "Record not found", http.StatusNotFound, CategoryNotFound)
	}

	// MongoDB: write errors (duplicate key, etc.)
	var writeErr mongo.WriteException
	if errors.As(err, &writeErr) {
		return parseMongoWriteError(writeErr)
	}

	// MongoDB: command errors
	var cmdErr mongo.CommandError
	if errors.As(err, &cmdErr) {
		return parseMongoCommandError(cmdErr)
	}

	// Default database error
	return NewAppError(err, "Database operation failed", http.StatusInternalServerError, CategoryDatabase)
}

// parseMongoWriteError handles MongoDB write errors
func parseMongoWriteError(err mongo.WriteException) *AppError {
	for _, we := range err.WriteErrors {
		switch we.Code {
		// Duplicate key error
		case 11000:
			return NewAppError(err, "Duplicate entry", http.StatusConflict, CategoryValidation)
		// Document validation failure
		case 121:
			return NewAppError(err, "Data validation failed", http.StatusBadRequest, CategoryValidation)
		}
	}

	return NewAppError(err, "Database write error", http.StatusInternalServerError, CategoryDatabase)
}

// parseMongoCommandError handles MongoDB command errors
func parseMongoCommandError(err mongo.CommandError) *AppError {
	switch {
	case err.HasErrorCode(13): // Unauthorized
		return NewAppError(err, "Database access denied", http.StatusForbidden, CategoryPermission)
	case err.HasErrorCode(50): // MaxTimeMSExpired
		return NewAppError(err, "Operation timed out", http.StatusGatewayTimeout, CategoryTimeout)
	default:
		return NewAppError(err, "Internal database error", http.StatusInternalServerError, CategoryDatabase)
	}
}

// parseConstraintError parses constraint-like errors from error messages
func parseConstraintError(err error) *AppError {
	errMsg := err.Error()

	switch {
	case strings.Contains(errMsg, "duplicate"):
		return NewAppError(err, "Record already exists", http.StatusConflict, CategoryValidation)
	case strings.Contains(errMsg, "validation"):
		return NewAppError(err, "Data validation failed", http.StatusBadRequest, CategoryValidation)
	default:
		return NewAppError(err, "Data constraint violation", http.StatusBadRequest, CategoryValidation)
	}
}

// IsDatabaseError checks if error is database related error
func IsDatabaseError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Category == CategoryDatabase
	}

	// Check for MongoDB error types
	var writeErr mongo.WriteException
	var cmdErr mongo.CommandError
	return errors.Is(err, mongo.ErrNoDocuments) ||
		errors.As(err, &writeErr) ||
		errors.As(err, &cmdErr)
}

// GetStatusCode returns proper HTTP status code for error
func GetStatusCode(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.StatusCode
	}

	if databaseErr := WrapDatabaseError(err); databaseErr != nil {
		return databaseErr.StatusCode
	}

	return http.StatusInternalServerError
}

// GetErrorMessage returns user friendly error message
func GetErrorMessage(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Message
	}

	if databaseErr := WrapDatabaseError(err); databaseErr != nil {
		return databaseErr.Message
	}

	return "Internal server error"
}

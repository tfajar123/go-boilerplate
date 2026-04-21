package utils

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"go-boilerplate/ent"

	"github.com/lib/pq"
)

// Database Error Types
var (
	ErrDatabaseConnection  = errors.New("failed to connect to database")
	ErrDatabaseQuery       = errors.New("database query error")
	ErrNotFound            = errors.New("record not found")
	ErrDuplicateEntry      = errors.New("duplicate entry")
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrConstraintViolation = errors.New("constraint violation")
	ErrInvalidInput        = errors.New("invalid input data")
	ErrTransactionFailed   = errors.New("transaction failed")
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

	// Ent ORM Errors
	if ent.IsNotFound(err) {
		return NewAppError(err, "Record not found", http.StatusNotFound, CategoryNotFound)
	}

	if ent.IsConstraintError(err) {
		return parseConstraintError(err)
	}

	if ent.IsValidationError(err) {
		return NewAppError(err, "Invalid data", http.StatusBadRequest, CategoryValidation)
	}

	// PostgreSQL Errors
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return parsePostgresError(pqErr)
	}

	// Standard SQL Errors
	if errors.Is(err, sql.ErrNoRows) {
		return NewAppError(err, "Record not found", http.StatusNotFound, CategoryNotFound)
	}

	if errors.Is(err, sql.ErrConnDone) {
		return NewAppError(err, "Database connection lost", http.StatusServiceUnavailable, CategoryDatabase)
	}

	if errors.Is(err, sql.ErrTxDone) {
		return NewAppError(err, "Transaction already completed", http.StatusInternalServerError, CategoryDatabase)
	}

	// Default database error
	return NewAppError(err, "Database operation failed", http.StatusInternalServerError, CategoryDatabase)
}

// parsePostgresError handles specific PostgreSQL error codes
func parsePostgresError(err *pq.Error) *AppError {
	switch err.Code {
	// Unique violation
	case "23505":
		return NewAppError(err, "Duplicate entry", http.StatusConflict, CategoryValidation)

	// Foreign key violation
	case "23503":
		return NewAppError(err, "Invalid reference data", http.StatusBadRequest, CategoryValidation)

	// Not null violation
	case "23502":
		return NewAppError(err, "Field cannot be empty", http.StatusBadRequest, CategoryValidation)

	// Check violation
	case "23514":
		return NewAppError(err, "Data validation failed", http.StatusBadRequest, CategoryValidation)

	// Invalid text representation
	case "22P02":
		return NewAppError(err, "Invalid data format", http.StatusBadRequest, CategoryValidation)

	// Connection errors
	case "08006", "08001", "08004", "57P01":
		return NewAppError(err, "Database connection failed", http.StatusServiceUnavailable, CategoryDatabase)

	// Timeout
	case "57014":
		return NewAppError(err, "Operation timed out", http.StatusGatewayTimeout, CategoryTimeout)

	// Permission denied
	case "42501":
		return NewAppError(err, "Database access denied", http.StatusForbidden, CategoryPermission)

	default:
		// Return generic database error
		return NewAppError(err, "Internal database error", http.StatusInternalServerError, CategoryDatabase)
	}
}

// parseConstraintError parses Ent constraint error
func parseConstraintError(err error) *AppError {
	errMsg := err.Error()

	switch {
	case strings.Contains(errMsg, "unique"):
		return NewAppError(err, "Record already exists", http.StatusConflict, CategoryValidation)
	case strings.Contains(errMsg, "foreign key"):
		return NewAppError(err, "Invalid data reference", http.StatusBadRequest, CategoryValidation)
	case strings.Contains(errMsg, "check constraint"):
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

	// Check for database error types
	_, isPQErr := err.(*pq.Error)
	return isPQErr ||
		ent.IsConstraintError(err) ||
		errors.Is(err, sql.ErrConnDone) ||
		errors.Is(err, sql.ErrNoRows) ||
		errors.Is(err, sql.ErrTxDone)
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

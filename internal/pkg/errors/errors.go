package errors

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// ErrorCode represents the type of error
type ErrorCode string

const (
	// Validation errors
	ValidationError     ErrorCode = "VALIDATION_ERROR"
	InvalidInput        ErrorCode = "INVALID_INPUT"
	MissingRequiredData ErrorCode = "MISSING_REQUIRED_DATA"
	InvalidFormat       ErrorCode = "INVALID_FORMAT"

	// Not found errors
	NotFoundError   ErrorCode = "NOT_FOUND"
	ProductNotFound ErrorCode = "PRODUCT_NOT_FOUND"
	UserNotFound    ErrorCode = "USER_NOT_FOUND"

	// Business logic errors
	BusinessLogicError ErrorCode = "BUSINESS_LOGIC_ERROR"
	InsufficientStock  ErrorCode = "INSUFFICIENT_STOCK"
	DuplicateEntry     ErrorCode = "DUPLICATE_ENTRY"

	// Database errors
	DatabaseError   ErrorCode = "DATABASE_ERROR"
	ConnectionError ErrorCode = "CONNECTION_ERROR"
	MigrationError  ErrorCode = "MIGRATION_ERROR"

	// Internal server errors
	InternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"
	UnknownError        ErrorCode = "UNKNOWN_ERROR"

	// Authentication/Authorization errors
	UnauthorizedError ErrorCode = "UNAUTHORIZED"
	ForbiddenError    ErrorCode = "FORBIDDEN"
	TokenExpiredError ErrorCode = "TOKEN_EXPIRED"
)

// AppError represents a custom application error
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	StatusCode int       `json:"-"`
	Details    any       `json:"details,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewAppError creates a new AppError
func NewAppError(code ErrorCode, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details any) *AppError {
	e.Details = details
	return e
}

// Validation Error Creators
func NewValidationError(message string) *AppError {
	return NewAppError(ValidationError, message, fiber.StatusBadRequest)
}

func NewInvalidInputError(message string) *AppError {
	return NewAppError(InvalidInput, message, fiber.StatusBadRequest)
}

func NewMissingRequiredDataError(field string) *AppError {
	return NewAppError(MissingRequiredData, fmt.Sprintf("Missing required field: %s", field), fiber.StatusBadRequest)
}

func NewInvalidFormatError(field string) *AppError {
	return NewAppError(InvalidFormat, fmt.Sprintf("Invalid format for field: %s", field), fiber.StatusBadRequest)
}

// Not Found Error Creators
func NewNotFoundError(message string) *AppError {
	return NewAppError(NotFoundError, message, fiber.StatusNotFound)
}

func NewProductNotFoundError(id string) *AppError {
	return NewAppError(ProductNotFound, fmt.Sprintf("Product with ID %s not found", id), fiber.StatusNotFound)
}

func NewUserNotFoundError(id string) *AppError {
	return NewAppError(UserNotFound, fmt.Sprintf("User with ID %s not found", id), fiber.StatusNotFound)
}

// Business Logic Error Creators
func NewBusinessLogicError(message string) *AppError {
	return NewAppError(BusinessLogicError, message, fiber.StatusUnprocessableEntity)
}

func NewInsufficientStockError(available, required int) *AppError {
	return NewAppError(InsufficientStock,
		fmt.Sprintf("Insufficient stock. Available: %d, Required: %d", available, required),
		fiber.StatusConflict).WithDetails(map[string]int{
		"available": available,
		"required":  required,
	})
}

func NewDuplicateEntryError(field, value string) *AppError {
	return NewAppError(DuplicateEntry,
		fmt.Sprintf("Duplicate entry for %s: %s", field, value),
		fiber.StatusConflict)
}

// Database Error Creators
func NewDatabaseError(message string) *AppError {
	return NewAppError(DatabaseError, message, fiber.StatusInternalServerError)
}

func NewConnectionError(message string) *AppError {
	return NewAppError(ConnectionError, message, fiber.StatusServiceUnavailable)
}

func NewMigrationError(message string) *AppError {
	return NewAppError(MigrationError, message, fiber.StatusInternalServerError)
}

// Internal Server Error Creators
func NewInternalServerError(message string) *AppError {
	return NewAppError(InternalServerError, message, fiber.StatusInternalServerError)
}

func NewUnknownError() *AppError {
	return NewAppError(UnknownError, "An unknown error occurred", fiber.StatusInternalServerError)
}

// Authentication/Authorization Error Creators
func NewUnauthorizedError(message string) *AppError {
	if message == "" {
		message = "Unauthorized access"
	}
	return NewAppError(UnauthorizedError, message, fiber.StatusUnauthorized)
}

func NewForbiddenError(message string) *AppError {
	if message == "" {
		message = "Access forbidden"
	}
	return NewAppError(ForbiddenError, message, fiber.StatusForbidden)
}

func NewTokenExpiredError() *AppError {
	return NewAppError(TokenExpiredError, "Token has expired", fiber.StatusUnauthorized)
}

// FromError converts a standard error to AppError
func FromError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return NewUnknownError().WithDetails(err.Error())
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetStatusCode returns the HTTP status code for the error
func GetStatusCode(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.StatusCode
	}
	return http.StatusInternalServerError
}

package errorx

import "errors"

// AppError represents a business logic error with a user-facing message.
type AppError struct {
	Code    string
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// Sentinel errors for common cases.
var (
	ErrNotFound   = New("NOT_FOUND", "Resource not found")
	ErrExpired    = New("EXPIRED", "Resource has expired")
	ErrExhausted  = New("EXHAUSTED", "Resource has reached its usage limit")
	ErrForbidden  = New("FORBIDDEN", "You do not have permission")
	ErrConflict   = New("CONFLICT", "Resource already exists")
	ErrBadRequest = New("BAD_REQUEST", "Invalid request")
	ErrInternal   = New("INTERNAL", "Internal server error")
)

// Is checks whether the target is an AppError with the same code.
func Is(err error, target *AppError) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == target.Code
	}
	return false
}

// Wrap creates a new AppError with a custom message but same code as the base.
func Wrap(base *AppError, message string) *AppError {
	return &AppError{Code: base.Code, Message: message}
}

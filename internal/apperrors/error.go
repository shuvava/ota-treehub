package apperrors

import (
	"fmt"

	"github.com/shuvava/treehub/internal/data"
)

// AppErrorCode is a type of AppError
// each code should be in format "namespace:error_code"
// namespace value should be unique across application
// error_code should be unique across package
type AppErrorCode string

// AppError application error with additional details
type AppError struct {
	// ErrorCode application error code
	ErrorCode AppErrorCode `json:"error_code"`
	// Description description of error
	Description string `json:"description"`
	// ErrorID unique ID of error required for easier look error in application logs
	ErrorID data.CorrelationID `json:"error_id"`
}

func (err AppError) Error() string {
	return fmt.Sprintf("(%s) : %s", err.ErrorCode, err.Description)
}

// NewAppError creates new AppError
func NewAppError(code AppErrorCode, descr string) error {
	id := data.NewCorrelationID()
	return AppError{
		ErrorCode:   code,
		Description: descr,
		ErrorID:     id,
	}
}

// CreateError create new AppError
func CreateError(code AppErrorCode, descr string, err error) error {
	return NewAppError(code, fmt.Sprintf("%s (%v)", descr, err))
}

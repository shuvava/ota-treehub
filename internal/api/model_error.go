package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shuvava/go-logging/logger"

	"github.com/shuvava/treehub/internal/apperrors"
	"github.com/shuvava/treehub/internal/data"
)

// ErrorResponse is http error response model
type ErrorResponse struct {
	// ErrorCode application error code
	ErrorCode string `json:"error_code"`
	// StatusCode HTTP response status code
	StatusCode int `json:"status_code"`
	// Description description of error
	Description string `json:"description"`
	// RequestID HTTP requestID go from header of request
	RequestID string `json:"request_id"`
	// ErrorID unique ID of error required for easier look error in application logs
	ErrorID string `json:"error_id"`
}

// NewErrorResponse creates new error response from error
func NewErrorResponse(ctx context.Context, statusCode int, err error) ErrorResponse {
	requestID := logger.GetRequestID(ctx)
	resp := ErrorResponse{
		StatusCode: statusCode,
		RequestID:  requestID,
	}

	var typedErr apperrors.AppError
	if errors.As(err, &typedErr) {
		resp.ErrorCode = string(typedErr.ErrorCode)
		resp.Description = typedErr.Description
		resp.ErrorID = typedErr.ErrorID.String()
	} else {
		resp.ErrorCode = apperrors.ErrorGeneric
		resp.Description = err.Error()
		resp.ErrorID = data.NewCorrelationID().String()
	}

	return resp
}

// EchoResponse build custom error response on err
func EchoResponse(ctx echo.Context, err error) error {
	c := GetRequestContext(ctx)
	var typedErr apperrors.AppError
	if !errors.As(err, &typedErr) {
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(c, http.StatusInternalServerError, err))
	}
	switch typedErr.ErrorCode {
	case apperrors.ErrorDataRefValidation, apperrors.ErrorDataObjectIDSerialization:
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(c, http.StatusBadRequest, err))
	default:
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(c, http.StatusInternalServerError, err))
	}
}

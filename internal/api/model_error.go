package api

import (
	"errors"
	"net/http"

	"github.com/shuvava/treehub/pkg/data"
	"github.com/shuvava/treehub/pkg/services"

	"github.com/labstack/echo/v4"

	cmnapi "github.com/shuvava/go-ota-svc-common/api"
	"github.com/shuvava/go-ota-svc-common/apperrors"
)

// EchoResponse build custom error response on err
func EchoResponse(ctx echo.Context, err error) error {
	c := cmnapi.GetRequestContext(ctx)
	var typedErr apperrors.AppError
	if !errors.As(err, &typedErr) {
		return ctx.JSON(http.StatusInternalServerError, cmnapi.NewErrorResponse(c, http.StatusInternalServerError, err))
	}
	switch typedErr.ErrorCode {
	case apperrors.ErrorDataValidation, apperrors.ErrorDataSerialization, data.ErrorDataSerializationObjectID, services.ErrorDataValidationRef:
		return ctx.JSON(http.StatusBadRequest, cmnapi.NewErrorResponse(c, http.StatusBadRequest, err))
	default:
		return ctx.JSON(http.StatusInternalServerError, cmnapi.NewErrorResponse(c, http.StatusInternalServerError, err))
	}
}

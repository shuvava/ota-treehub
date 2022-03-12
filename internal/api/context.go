package api

import (
	"fmt"
	"strconv"

	cmnapi "github.com/shuvava/go-ota-svc-common/api"

	"github.com/labstack/echo/v4"
	"github.com/shuvava/treehub/pkg/data"
)

const (
	headerForcePush = "x-ats-ostree-force"

	pathOPrefix = "oprefix"
	pathOSuffix = "osuffix"
)

// GetObjectID builds data.ObjectID from request path
func GetObjectID(ctx echo.Context) (data.ObjectID, error) {
	oprefix := ctx.Param(pathOPrefix)
	osuffix := ctx.Param(pathOSuffix)
	return data.NewObjectID(oprefix + osuffix)
}

// ValidateUploadContentType that Request has valid ContentType
func ValidateUploadContentType(ctx echo.Context) error {
	mime := cmnapi.GetContentType(ctx)
	if mime != echo.MIMEOctetStream {
		return fmt.Errorf("header %s mast be '%s' type", echo.HeaderContentType, echo.MIMEOctetStream)
	}
	return nil
}

// IsForcePush check if force push header was set in request
func IsForcePush(ctx echo.Context) bool {
	val := ctx.Request().Header.Get(headerForcePush)
	res, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}
	return res
}

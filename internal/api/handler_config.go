package api

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	cmnapi "github.com/shuvava/go-ota-svc-common/api"
)

const (
	// PathConfig is the path to the config file
	PathConfig = "/config"
	configFile = `[core]
repo_version=1
mode=archive-z2
`
)

// ConfigDownload is endpoint download OSTree config file from server to client
func ConfigDownload(ctx echo.Context) error {
	c := cmnapi.GetRequestContext(ctx)
	reader := strings.NewReader(configFile)
	if _, err := reader.WriteTo(ctx.Response().Writer); err != nil {
		return ctx.JSON(http.StatusInternalServerError, cmnapi.NewErrorResponse(c, http.StatusInternalServerError, err))
	}
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEOctetStream)
	ctx.Response().WriteHeader(http.StatusOK)
	ctx.Response().Flush()
	return nil
}

package api

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
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
	c := GetRequestContext(ctx)
	reader := strings.NewReader(configFile)
	if _, err := reader.WriteTo(ctx.Response().Writer); err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(c, http.StatusInternalServerError, err))
	}
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEOctetStream)
	ctx.Response().WriteHeader(http.StatusOK)
	ctx.Response().Flush()
	return nil
}

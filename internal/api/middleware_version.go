package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/shuvava/treehub/pkg/version"
)

// ServerHeader middleware adds a `Server` header to the response.
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	appVer := fmt.Sprintf("%s/%s", version.AppName, version.Version)
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, appVer)
		return next(c)
	}
}

package app

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo-contrib/prometheus"

	"github.com/shuvava/treehub/internal/api"
)

const (
	routeApiVer2 = "/api/v2"
	routeApiVer3 = "/api/v3"
)

// initWebServer creates echo http server and set request handlers
func (s *Server) initWebServer() {
	// Initialize Echo, set error handler, add in middleware
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = api.NewErrorHandler().Handler
	e.Pre(middleware.RemoveTrailingSlash())
	// logger Middleware (https://echo.labstack.com/middleware/logger/)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Gzip())
	// Server header
	e.Use(api.ServerHeader)

	initHealthRoutes(s, e)
	v2Group := e.Group(routeApiVer2, middleware.RequestID())
	initObjectRoutes(s, v2Group)
	initRefsRoutes(s, v2Group)
	initConfRoutes(v2Group)
	v3Group := e.Group(routeApiVer3, middleware.RequestID())
	initObjectRoutes(s, v3Group)
	initRefsRoutes(s, v3Group)
	initConfRoutes(v3Group)

	// Enable metrics middleware
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	s.Echo = e
}

func initHealthRoutes(s *Server, e *echo.Echo) {
	// Define a separate root 'health' group without the logging middleware added (for healthz/readyz)
	healthGroup := e.Group("")
	healthGroup.GET(api.LivenessPath, api.HealthzHandler)
	healthGroup.GET(api.ReadinessPath, api.ReadyzHandler(
		func(ctx context.Context) api.HealthEntryStatus {
			resource := "repository"
			if err := s.svc.Db.Ping(ctx); err != nil {
				return api.HealthEntryStatus{
					Status:   api.StatusUnhealthy,
					Data:     err.Error(),
					Resource: resource,
				}
			}
			return api.HealthEntryStatus{
				Status:   api.StatusHealthy,
				Resource: resource,
			}
		},
	))
}

func initObjectRoutes(s *Server, group *echo.Group) {
	group.GET(api.PathObject, func(c echo.Context) error {
		return api.ObjectDownload(c, s.svc.Objects)
	})
	group.POST(api.PathObject, func(c echo.Context) error {
		return api.ObjectUpload(c, s.svc.Objects)
	})
	group.PUT(api.PathObject, func(c echo.Context) error {
		return api.ObjectUploadCompleted(c, s.svc.Objects)
	})
	group.HEAD(api.PathObject, func(c echo.Context) error {
		return api.ObjectExists(c, s.svc.Objects)
	})
}

func initRefsRoutes(s *Server, group *echo.Group) {
	group.POST(api.PathRefs, func(c echo.Context) error {
		return api.RefsUpload(c, s.svc.Refs)
	})
	group.GET(api.PathRefs, func(c echo.Context) error {
		return api.RefDownload(c, s.svc.Refs)
	})
}

func initConfRoutes(group *echo.Group) {
	group.GET(api.PathConfig, func(c echo.Context) error {
		return api.ConfigDownload(c)
	})
}

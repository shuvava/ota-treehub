package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/shuvava/treehub/pkg/services"

	"github.com/shuvava/treehub/internal/blobs"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"

	"github.com/shuvava/treehub/internal/config"
	intDb "github.com/shuvava/treehub/internal/db"
	"github.com/shuvava/treehub/internal/logger"
)

// Server is main application servers
type Server struct {
	Echo   *echo.Echo
	log    logger.Logger
	config *config.AppConfig
	mu     sync.Mutex
	svc    struct {
		Db          intDb.BaseRepository
		ObjectRepo  intDb.ObjectRepository
		RefRepo     intDb.RefRepository
		ObjectStore blobs.ObjectStore
		Objects     *services.ObjectService
		Refs        *services.RefService
	}
}

// NewServer creates new Server instance
func NewServer(logger logger.Logger) *Server {
	s := &Server{
		log: logger.SetContext("internal-http"),
	}

	s.initWebServer()
	s.initConfig()

	return s
}

// initConfig load app config
func (s *Server) initConfig() {
	cfg := config.NewConfig(s.log, s.OnConfigChange)
	s.OnConfigChange(cfg)
}

// OnConfigChange execute operation required on config change
func (s *Server) OnConfigChange(newCfg *config.AppConfig) {
	s.mu.Lock()
	lvl := logger.ToLogLevel(newCfg.LogLevel)
	_ = s.log.SetLevel(lvl)
	s.config = newCfg
	s.initServices()
	s.mu.Unlock()
}

// Start starts web server main event loop
func (s *Server) Start() {
	// Determine API listen address/port
	serverListenAddr := fmt.Sprintf("0.0.0.0:%d", s.config.Port)
	// Start server
	go func() {
		if err := s.Echo.Start(serverListenAddr); err != nil {
			s.log.WithError(err).
				Fatal("Fatal error in API server")
		}
	}()
	logrus.Info(fmt.Sprintf("Service start listening on %s", serverListenAddr))
	// Wait for interrupt signal to gracefully shutting down the web server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancelShutdown := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelShutdown()
	if err := s.Echo.Shutdown(ctx); err != nil {
		s.log.WithError(err).
			Fatal("Error shutting down API server")
	}
}

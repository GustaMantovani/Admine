package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/GustaMantovani/Admine/server_handler/internal/config"
	"github.com/GustaMantovani/Admine/server_handler/pkg"
)

// Server represents the HTTP API server
type Server struct {
	server *http.Server
}

// NewServer creates a new API server instance
func NewServer(cfg *config.Config) *Server {
	router := SetupRoutes()

	address := fmt.Sprintf("%s:%d", cfg.WebSever.Host, cfg.WebSever.Port)

	server := &http.Server{
		Addr:    address,
		Handler: router,
		// Timeouts
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		server: server,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	pkg.Logger.Info("Starting HTTP server at %s", s.server.Addr)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// StartBackground starts the HTTP server in the background
func (s *Server) StartBackground() error {
	pkg.Logger.Info("Starting HTTP server in background at %s", s.server.Addr)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			pkg.Logger.Error("Server error: %v", err)
		}
	}()

	return nil
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	pkg.Logger.Info("Stopping HTTP server")

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	pkg.Logger.Info("HTTP server stopped")
	return nil
}

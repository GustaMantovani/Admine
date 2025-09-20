package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/GustaMantovani/Admine/server_handler/internal"
)

// Server represents the HTTP API server
type Server struct {
	server *http.Server
}

// NewServer creates a new API server instance
func NewServer() *Server {
	cfg := internal.Get().Config
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
	slog.Info("Starting HTTP server", "addr", s.server.Addr)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// StartBackground starts the HTTP server in the background
func (s *Server) StartBackground() error {
	slog.Info("Starting HTTP server in background", "addr", s.server.Addr)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server error", "error", err)
		}
	}()

	return nil
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	slog.Info("Stopping HTTP server")

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	slog.Info("HTTP server stopped")
	return nil
}

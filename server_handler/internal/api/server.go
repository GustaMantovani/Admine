package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/GustaMantovani/Admine/server_handler/internal/config"
)

// Server wraps the HTTP server
type Server struct {
	server *http.Server
}

// NewServer creates a new API Server
func NewServer(cfg config.WebServerConfig, router http.Handler) *Server {
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	s := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{server: s}
}

// Start starts the HTTP server (blocking)
func (s *Server) Start() error {
	slog.Info("Starting HTTP server", "addr", s.server.Addr)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// StartBackground starts the HTTP server in a goroutine
func (s *Server) StartBackground() error {
	slog.Info("Starting HTTP server in background", "addr", s.server.Addr)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server error", "error", err)
		}
	}()

	return nil
}

// Stop gracefully shuts down the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	slog.Info("Stopping HTTP server")

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	slog.Info("HTTP server stopped")
	return nil
}

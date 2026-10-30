// Package server runs an HTTP handler and drains it on shutdown.
package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/LoyalMinion/api-streaming-availability/internal/config"
)

// Server serves a handler until its context ends, then drains in-flight requests.
type Server struct {
	logger *slog.Logger
	server *http.Server
	drain  time.Duration
}

// New returns a Server for handler, configured by cfg.
func New(cfg config.Config, handler http.Handler, logger *slog.Logger) *Server {
	return &Server{
		logger: logger,
		server: &http.Server{
			Addr:              ":" + cfg.Port,
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
		},
		drain: cfg.DrainTimeout,
	}
}

// ListenAndServe listens on the configured port and serves until ctx ends,
// then waits up to the drain timeout for in-flight requests to finish.
func (s *Server) ListenAndServe(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	return s.Serve(ctx, listener)
}

// Serve serves listener until ctx ends, then waits up to the drain timeout
// for in-flight requests to finish.
func (s *Server) Serve(ctx context.Context, listener net.Listener) error {
	s.logger.Info("listen", "addr", listener.Addr().String())

	served := make(chan error, 1)
	go func() {
		served <- s.server.Serve(listener)
	}()

	select {
	case err := <-served:
		return fmt.Errorf("serve: %w", err)
	case <-ctx.Done():
	}

	s.logger.Info("drain", "timeout", s.drain.String())

	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), s.drain)
	defer cancel()

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	return nil
}

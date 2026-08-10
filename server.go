package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	httpServer     *http.Server
	logger         *slog.Logger
	shutdownPeriod func(context.Context) (context.Context, context.CancelFunc)
}

func New(handler http.Handler, address string, shutdownTimeout time.Duration, logger *slog.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              address,
			Handler:           handler,
			ReadHeaderTimeout: 5 * shutdownTimeout / 2,
		},
		logger: logger,
		shutdownPeriod: func(ctx context.Context) (context.Context, context.CancelFunc) {
			return context.WithTimeout(ctx, shutdownTimeout)
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	errorsCh := make(chan error, 1)
	go func() {
		s.logger.Info("server started", "address", s.httpServer.Addr)
		errorsCh <- s.httpServer.ListenAndServe()
	}()

	select {
	case err := <-errorsCh:
		return err
	case <-ctx.Done():
		s.logger.Info("server shutting down")
	}

	shutdownCtx, cancel := s.shutdownPeriod(context.Background())
	defer cancel()
	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		return errors.Join(http.ErrServerClosed, err)
	}

	return http.ErrServerClosed
}

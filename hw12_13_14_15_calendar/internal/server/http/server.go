package internalhttp

import (
	"context"
	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"log/slog"
	"net"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	logger     slog.Logger
}

func NewServer(logger slog.Logger, cfg config.Config) *Server {
	return &Server{
		logger: logger,
		httpServer: &http.Server{
			Addr:         net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
			ReadTimeout:  cfg.Server.Timeout,
			WriteTimeout: cfg.Server.Timeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
	}
}

func (s *Server) Start(mux *runtime.ServeMux) error {
	s.httpServer.Handler = loggingMiddleware(mux)
	s.logger.Info("starting http server with address", "address", s.httpServer.Addr)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error("could not listen on", "address", s.httpServer.Addr, ":", err)
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("shutting down http server")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("failed to shutdown http server:", err)
	}
	s.logger.Info("shutting down http server gracefully")
	return nil
}

package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"proxy/internal/config"
	"proxy/internal/logger"
	"proxy/internal/proxy"
)

type Server struct {
	addr    string
	handler http.Handler
	srv     *http.Server
}

func NewServer(cfg *config.Config) (*Server, error) {
	proxyHandler, err := proxy.NewProxyHandler(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create proxy handler: %v", err)
	}

	handler := loggingMiddleware(proxyHandler)

	return &Server{
		addr:    cfg.HTTPServer.Address,
		handler: handler,
		srv:     &http.Server{Addr: cfg.HTTPServer.Address, Handler: handler},
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	l := logger.LoggerFromContext(ctx)
	go func() {
		l.Info(fmt.Sprintf("Starting server on %s", s.addr))
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Info(fmt.Sprintf("server startup failed: %v", err))
			os.Exit(1)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		return s.shutdown()
	case sig := <-stopChan:
		l.Info(fmt.Sprintf("received signal %s", sig))
		return s.shutdown()
	}
}

func (s *Server) shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server shutdown failed: %v", err)
	}
	return nil
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Logger().Info(
			"Received request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_ip", proxy.ReadUserIP(r)),
		)
		next.ServeHTTP(w, r)
	})
}

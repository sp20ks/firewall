// server/server.go
package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"firewall/internal/config"
	"firewall/internal/proxy"
)

type Server struct {
	addr    string
	handler http.Handler
	srv     *http.Server
}

func NewServer(cfg *config.Config) (*Server, error) {
	proxyResources := make([]proxy.Resource, len(cfg.Resources))
	for i, r := range cfg.Resources {
		proxyResources[i] = proxy.Resource{
			Host:     r.Host,
			Endpoint: r.Endpoint,
		}
	}

	proxyHandler, err := proxy.NewProxyHandler(proxyResources)
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
	go func() {
		log.Printf("Starting server on %s", s.addr)
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server startup failed: %v", err)
			os.Exit(1)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		return s.shutdown()
	case sig := <-stopChan:
		log.Printf("received signal %s", sig)
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
	log.Println("server stopped")
	return nil
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: method=%s, path=%s, remote=%s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

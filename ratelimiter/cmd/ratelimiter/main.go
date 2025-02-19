package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"ratelimiter/internal/config"
	"ratelimiter/internal/ratelimiter"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	ipRateLimiter := ratelimiter.NewIPRateLimiter(cfg.RedisAddr)
	defer func() {
		if err := ipRateLimiter.Close(); err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/rate_limit", ratelimiter.HandleCheckLimit(ipRateLimiter))

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: mux,
	}

	go func() {
		log.Printf("Starting server on %s", cfg.Address)	
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Failed to shutdown server: %v", err)
	}
}

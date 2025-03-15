package main

import (
	"auth/internal/auth"
	"auth/internal/config"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	a := auth.NewAuth(cfg.SecretKey, time.Duration(cfg.KeyTTL))
	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth", auth.HandleGetJwtToken(a))
	mux.HandleFunc("GET /verify", auth.VerifyJwtToken(a))

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: mux,
	}

	go func() {
		log.Printf("Starting server on %s", cfg.Address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
			os.Exit(1)
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

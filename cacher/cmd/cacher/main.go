package main

import (
	"cacher/internal/cacher"
	"cacher/internal/config"
	"cacher/internal/logger"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	redisCacher, err := cacher.NewCacher(cfg.RedisAddr, time.Duration(cfg.Ttl)*time.Second)
	if err != nil {
		log.Fatalf("Error initializing cacher: %v", err)
	}

	defer func() {
		if err := redisCacher.Close(); err != nil {
			logger.Logger().Info("error closing Redis connection", zap.Error(err))
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /cache", cacher.HandleGetCache(redisCacher))
	mux.HandleFunc("POST /cache", cacher.HandleSetCache(redisCacher))

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: mux,
	}

	go func() {
		logger.Logger().Info(fmt.Sprintf("starting server on %s", cfg.Address))
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
		logger.Logger().Info("failed to shutdown server", zap.Error(err))
	}
}

package main

import (
	"auth/internal/config"
	"auth/internal/delivery"
	"auth/internal/logger"
	"auth/internal/repository/postgres"
	"auth/internal/usecase"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.AuthDB.Host, cfg.AuthDB.Port, cfg.AuthDB.User, cfg.AuthDB.Password, cfg.AuthDB.DBName, cfg.AuthDB.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	userRepo := postgres.NewPostgresUserRepository(db)
	authUseCase := usecase.NewAuthUseCase(userRepo, cfg.SecretKey, time.Duration(cfg.KeyTTL))
	authHandler := delivery.NewAuthHandler(authUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth", authHandler.HandleGetJwtToken)
	mux.HandleFunc("POST /register", authHandler.HandleRegisterUser)
	mux.HandleFunc("GET /verify", authHandler.VerifyJwtToken)

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

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	authservice "rules-engine/internal/clients/auth_service"
	"rules-engine/internal/config"
	"rules-engine/internal/delivery"
	"rules-engine/internal/delivery/middleware"
	"rules-engine/internal/repository/postgres"
	"rules-engine/internal/usecase"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	resourceRepo := postgres.NewPostgresResourceRepository(db)
	resourceUseCase := usecase.NewResorceUseCase(resourceRepo)
	resourceHandler := delivery.NewResourceHandler(resourceUseCase)

	ipListRepo := postgres.NewPostgresIPListRepository(db)
	ipListUseCase := usecase.NewIPListUseCase(ipListRepo)
	ipListHandler := delivery.NewIPListHandler(ipListUseCase)

	authClient := authservice.NewAuthClient(cfg.AuthURL)
	authMiddleware := middleware.AuthMiddleware(authClient)

	mux := http.NewServeMux()
	mux.Handle("POST /resources", authMiddleware(http.HandlerFunc(resourceHandler.HandleCreateResource)))
	mux.Handle("PUT /resources/{id}", authMiddleware(http.HandlerFunc(resourceHandler.HandleUpdateResource)))
	mux.HandleFunc("GET /resources", resourceHandler.HandleGetActiveResources)
	mux.Handle("POST /ip_lists", authMiddleware(http.HandlerFunc(ipListHandler.HandleCreateIPList)))
	mux.Handle("PUT /ip_lists/{id}", authMiddleware(http.HandlerFunc(ipListHandler.HandleUpdateIPList)))
	mux.HandleFunc("GET /ip_lists", ipListHandler.HandleGetIPLists)

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

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
	"rules-engine/internal/logger"
	"rules-engine/internal/repository/postgres"
	"rules-engine/internal/usecase"
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
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	resourceRepo := postgres.NewPostgresResourceRepository(db)
	ipListRepo := postgres.NewPostgresIPListRepository(db)
	resourceIPListRepo := postgres.NewPostgresResourceIPListRepository(db)
	resourceRuleRepo := postgres.NewPostgresResourceRuleRepository(db)
	ruleRepo := postgres.NewPostgresRuleRepository(db)

	ipListUseCase := usecase.NewIPListUseCase(ipListRepo)
	ruleUseCase := usecase.NewRuleUseCase(ruleRepo)
	resourceUseCase := usecase.NewResourceUseCase(resourceRepo, ipListUseCase, ruleUseCase, resourceIPListRepo, resourceRuleRepo)
	analyzer := usecase.NewAnalyzerUseCase(ruleRepo, ipListRepo)

	resourceHandler := delivery.NewResourceHandler(resourceUseCase)
	ipListHandler := delivery.NewIPListHandler(ipListUseCase)
	ruleHandler := delivery.NewRuleHandler(ruleUseCase)
	analyzerHandler := delivery.NewAnalyzerHandler(analyzer)

	authClient := authservice.NewAuthClient(cfg.AuthURL)
	authMiddleware := middleware.AuthMiddleware(authClient)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /resources", resourceHandler.HandleGetActiveResources)
	mux.HandleFunc("GET /resources/{id}", resourceHandler.HandleGetResourceByID)
	mux.Handle("POST /resources", authMiddleware(http.HandlerFunc(resourceHandler.HandleCreateResource)))
	mux.Handle("POST /resources/{id}/attach_ip_list", authMiddleware(http.HandlerFunc(resourceHandler.HandleAttachIPList)))
	mux.Handle("POST /resources/{id}/detach_ip_list", authMiddleware(http.HandlerFunc(resourceHandler.HandleDetachIPList)))
	mux.Handle("POST /resources/{id}/attach_rule", authMiddleware(http.HandlerFunc(resourceHandler.HandleAttachRule)))
	mux.Handle("POST /resources/{id}/detach_rule", authMiddleware(http.HandlerFunc(resourceHandler.HandleDetachRule)))
	mux.Handle("PUT /resources/{id}", authMiddleware(http.HandlerFunc(resourceHandler.HandleUpdateResource)))

	mux.Handle("POST /ip_lists", authMiddleware(http.HandlerFunc(ipListHandler.HandleCreateIPList)))
	mux.Handle("PUT /ip_lists/{id}", authMiddleware(http.HandlerFunc(ipListHandler.HandleUpdateIPList)))
	mux.HandleFunc("GET /ip_lists", ipListHandler.HandleGetIPLists)

	mux.Handle("POST /rules", authMiddleware(http.HandlerFunc(ruleHandler.HandleCreateRule)))
	mux.Handle("PUT /rules/{id}", authMiddleware(http.HandlerFunc(ruleHandler.HandleUpdateRule)))
	mux.HandleFunc("GET /rules", ruleHandler.HandleGetRules)

	mux.HandleFunc("GET /analyze", analyzerHandler.HandleAnalyzeRequest)

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: middleware.CorsMiddleware(mux),
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

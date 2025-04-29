package main

import (
	"context"
	"log"
	"proxy/internal/config"
	"proxy/internal/logger"
	"proxy/internal/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	server, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Error initializing server: %v", err)
	}

	l := logger.Logger()
	if err != nil {
		log.Fatalf("Error logger server: %v", err)
	}

	ctx := logger.ContextWithLogger(context.Background(), l)
	err = server.Start(ctx)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

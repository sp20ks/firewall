package main

import (
	"context"
	"log"
	"proxy/internal/config"
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

	err = server.Start(context.Background())
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

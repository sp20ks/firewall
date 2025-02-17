package main

import (
	"context"
	"proxy/internal/config"
	"proxy/internal/server"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading allowed IPs: %v", err)
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

package main

import (
	"firewall/internal/config"
	"firewall/internal/server"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading allowed IPs: %v", err)
	}

	err = server.Start(cfg)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

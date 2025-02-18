package main

import (
	"log"
	"net/http"
	"ratelimiter/internal/config"
	"ratelimiter/internal/ratelimiter"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	ipRateLimiter := ratelimiter.NewIPRateLimiter(cfg.MaxTokens, cfg.RefillRate)

	mux := http.NewServeMux()
	mux.HandleFunc("/rate_limit", ratelimiter.HandleCheckLimit(ipRateLimiter))

		log.Printf("Starting server on %s", cfg.Address)
		http.ListenAndServe(cfg.Address, mux)
		// if err != nil {
		// 	log.Fatalf("Error starting server: %v", err)
		// }
}

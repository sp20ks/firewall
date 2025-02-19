package ratelimiter

import (
	"log"
	"net/http"
)

func HandleCheckLimit(ipRateLimiter *IPRateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.URL.Query().Get("ip")
		if ip == "" {
			log.Printf("Missing required 'ip' parameter")
			http.Error(w, "Missing required 'ip' parameter", http.StatusBadRequest)
			return
		}

		if err := ipRateLimiter.Allow(ip); err != nil {
			log.Printf("Rate limit exceeded for IP %s: %v", ip, err)
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

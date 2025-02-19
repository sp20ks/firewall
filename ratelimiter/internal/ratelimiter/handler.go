package ratelimiter

import (
	"log"
	"net/http"
)

func HandleCheckLimit(ipRateLimiter *IPRateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.URL.Query().Get("ip")
		if ip == "" {
			log.Printf("Invalid IP received")
			http.Error(w, "Invalid IP", http.StatusInternalServerError)
			return
		}

		if err := ipRateLimiter.Allow(ip); err != nil {
			log.Printf("Rate Limit Exceeded IP=%s", ip)
			http.Error(w, "Rate Limit Exceeded", http.StatusTooManyRequests)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

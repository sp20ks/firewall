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

		if ipRateLimiter.Allow(ip) {
			w.WriteHeader(http.StatusOK)
		} else {
			log.Printf("Rate Limit Exceeded IP=%s", ip)
			http.Error(w, "Rate Limit Exceeded", http.StatusTooManyRequests)
		}
	}
}

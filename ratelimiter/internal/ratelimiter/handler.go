package ratelimiter

import (
	"net/http"
	"ratelimiter/internal/logger"

	"go.uber.org/zap"
)

func HandleCheckLimit(ipRateLimiter *IPRateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.URL.Query().Get("ip")
		l := logger.Logger()

		if ip == "" {
			l.Info("missing required 'ip' parameter")
			http.Error(w, "Missing required 'ip' parameter", http.StatusBadRequest)
			return
		}

		if err := ipRateLimiter.Allow(r.Context(), ip); err != nil {
			l.Info("rate limit exceeded for IP", zap.String("ip", ip), zap.Error(err))
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

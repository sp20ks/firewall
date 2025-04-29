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
			WriteJSONResponse(w, NewErrorResponse("missing required 'ip' parameter", http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := ipRateLimiter.Allow(r.Context(), ip); err != nil {
			l.Info("rate limit exceeded for IP", zap.String("ip", ip), zap.Error(err))
			WriteJSONResponse(w, NewErrorResponse("rate limit exceeded", http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		WriteJSONResponse(w, NewSuccessResponse("request allowed", http.StatusOK), http.StatusOK)
	}
}

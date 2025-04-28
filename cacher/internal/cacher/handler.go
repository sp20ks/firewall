package cacher

import (
	"cacher/internal/logger"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type CacheRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func HandleGetCache(c *Cacher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.Logger()
		key := r.URL.Query().Get("key")

		if key == "" {
			l.Info("missing required 'key' parameter")
			http.Error(w, "Missing required 'key' parameter", http.StatusBadRequest)
			return
		}

		val, err := c.GetCache(key)
		if err != nil {
			l.Info("cache miss for key", zap.String("key", key))
			http.Error(w, "Cache miss", http.StatusNoContent)
			return
		}

		l.Info("successfully handle getting cache by key", zap.String("key", key))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"key": key, "value": *val})
	}
}

func HandleSetCache(c *Cacher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.Logger()
		var req CacheRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			l.Info("Invalid request payload")
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if req.Key == "" || req.Value == "" {
			l.Info("Both 'key' and 'value' must be provided")
			http.Error(w, "Both 'key' and 'value' must be provided", http.StatusBadRequest)
			return
		}

		if err := c.SetCache(req.Key, req.Value); err != nil {
			l.Info("Failed to set cache for key", zap.String("key", req.Key), zap.Error(err))
			http.Error(w, "Failed to set cache", http.StatusInternalServerError)
			return
		}

		l.Info("successfully handle setting cache by key", zap.String("key", req.Key))

		w.WriteHeader(http.StatusCreated)
	}
}

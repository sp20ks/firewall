package cacher

import (
	"encoding/json"
	"log"
	"net/http"
)

type CacheRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func HandleGetCache(c *Cacher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if key == "" {
			log.Println("Missing required 'key' parameter")
			http.Error(w, "Missing required 'key' parameter", http.StatusBadRequest)
			return
		}

		val, err := c.GetCache(key)
		if err != nil {
			log.Printf("Cache miss for key %s", key)
			http.Error(w, "Cache miss", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"key": key, "value": *val})
	}
}

func HandleSetCache(c *Cacher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CacheRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("Invalid request payload")
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if req.Key == "" || req.Value == "" {
			log.Println("Both 'key' and 'value' must be provided")
			http.Error(w, "Both 'key' and 'value' must be provided", http.StatusBadRequest)
			return
		}

		if err := c.SetCache(req.Key, req.Value); err != nil {
			log.Printf("Failed to set cache for key %s: %v", req.Key, err)
			http.Error(w, "Failed to set cache", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

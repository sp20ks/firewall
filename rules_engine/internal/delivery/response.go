package delivery

import (
	"encoding/json"
	"net/http"
	"rules-engine/internal/logger"

	"go.uber.org/zap"
)

type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func JSONResponse[T any](w http.ResponseWriter, status int, data T, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err != nil {
		logger.Logger().Info("http error", zap.Error(err))
		resp := APIResponse[T]{Success: false, Error: err.Error()}
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	resp := APIResponse[T]{Success: true, Data: data}
	_ = json.NewEncoder(w).Encode(resp)
}

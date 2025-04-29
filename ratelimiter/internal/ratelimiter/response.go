package ratelimiter

import (
	"encoding/json"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code"`
	Timestamp  string `json:"timestamp"`
}

type SuccessResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Timestamp  string `json:"timestamp"`
}

func WriteJSONResponse(w http.ResponseWriter, payload any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func NewErrorResponse(msg string, status int) ErrorResponse {
	return ErrorResponse{
		Error:      msg,
		StatusCode: status,
		Timestamp:  time.Now().Format(time.RFC3339),
	}
}

func NewSuccessResponse(msg string, status int) SuccessResponse {
	return SuccessResponse{
		Message:    msg,
		StatusCode: status,
		Timestamp:  time.Now().Format(time.RFC3339),
	}
}

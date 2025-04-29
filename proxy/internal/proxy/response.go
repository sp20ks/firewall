package proxy

import (
	"encoding/json"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code"`
	Timestamp  string `json:"timestamp"`
	RequestID  string `json:"request_id"`
}

type SuccessResponse struct {
	Data       any    `json:"data"`
	StatusCode int    `json:"status_code"`
	Timestamp  string `json:"timestamp"`
	RequestID  string `json:"request_id"`
}

func WriteJSONResponse(w http.ResponseWriter, payload any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func NewErrorResponse(err string, status int, requestID string) ErrorResponse {
	return ErrorResponse{
		Error:      err,
		StatusCode: status,
		Timestamp:  time.Now().Format(time.RFC3339),
		RequestID:  requestID,
	}
}

package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type ErrorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code"`
	Timestamp  string `json:"timestamp"`
	RequestID  string `json:"request_id,omitempty"`
}

func (ph *ProxyHandler) modifyRequest(ctx context.Context, r *http.Request, resource *Resource) (*http.Request, error) {
	url, err := url.Parse(resource.Host)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, r.Method, url.String(), r.Body)
	if err != nil {
		return nil, err
	}

	for k, v := range r.Header {
		req.Header[k] = v
	}

	req.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	req.Header.Set("X-Forwarded-Proto", r.Proto)

	return req, nil
}

func (ph *ProxyHandler) forwardRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	resp, err := ph.transport.RoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("proxy request failed: %v", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	return resp, nil
}

func copyResponse(w http.ResponseWriter, resp *http.Response) {
	for k, v := range resp.Header {
		w.Header()[k] = v
	}

	w.WriteHeader(resp.StatusCode)

	_, err := io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("failed to copy response body: %v", err)
	}
}

func WriteJSONResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")

	if status >= http.StatusBadRequest {
		log.Printf("Error response: status=%d, error=%+v", status, data)
	}

	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (ph *ProxyHandler) validateRequest(r *http.Request) (code int, err error) {
	ip := ReadUserIP(r)

	if allowed, err := ph.rateLimiterClient.CheckLimit(ip); !allowed {
		log.Printf("Rate limit error ip=%s: %v", ip, err)
		return http.StatusTooManyRequests, fmt.Errorf("too many requests")
	}

	return http.StatusOK, nil
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

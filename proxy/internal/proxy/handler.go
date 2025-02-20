package proxy

import (
	"context"
	"net/http"
	"proxy/internal/config"
	"time"

	ratelimiterservice "proxy/internal/clients/ratelimiter_service"

	"github.com/google/uuid"
)

type ProxyHandler struct {
	configs           map[string]*Resource
	transport         http.RoundTripper
	rateLimiterClient *ratelimiterservice.RateLimiterClient
}

func NewProxyHandler(resources []Resource, cfg *config.Config) (*ProxyHandler, error) {
	configs, err := loadConfigs(resources)
	if err != nil {
		return nil, err
	}

	return &ProxyHandler{
		configs: configs,
		transport: &http.Transport{
			DisableKeepAlives: true,
		},
		rateLimiterClient: ratelimiterservice.NewRateLimiterClient(cfg.RateLimiterURL),
	}, nil
}

func (ph *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := uuid.NewString()
	ctx = context.WithValue(ctx, "request-id", requestID)

	resource, ok := ph.configs[r.URL.Path]
	if !ok {
		errResp := ErrorResponse{
			Error:      "endpoint not found",
			StatusCode: http.StatusNotFound,
			Timestamp:  time.Now().Format(time.RFC3339),
			RequestID:  requestID,
		}
		WriteJSONResponse(w, errResp, http.StatusNotFound)
		return
	}

	if code, err := ph.validateRequest(r); err != nil {
		errResp := ErrorResponse{
			Error:      err.Error(),
			StatusCode: code,
			Timestamp:  time.Now().Format(time.RFC3339),
			RequestID:  requestID,
		}
		WriteJSONResponse(w, errResp, code)
		return
	}

	req, err := ph.modifyRequest(ctx, r, resource)
	if err != nil {
		errResp := ErrorResponse{
			Error:      "internal server error",
			StatusCode: http.StatusInternalServerError,
			Timestamp:  time.Now().Format(time.RFC3339),
			RequestID:  requestID,
		}
		WriteJSONResponse(w, errResp, http.StatusInternalServerError)
		return
	}

	resp, err := ph.forwardRequest(ctx, req)
	if err != nil {
		errResp := ErrorResponse{
			Error:      "proxy error",
			StatusCode: http.StatusInternalServerError,
			Timestamp:  time.Now().Format(time.RFC3339),
			RequestID:  requestID,
		}
		WriteJSONResponse(w, errResp, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	copyResponse(w, resp)
}

package proxy

import (
	"context"
	"io"
	"log"
	"net/http"
	"proxy/internal/config"
	"time"

	cacherservice "proxy/internal/clients/cacher_service"
	ratelimiterservice "proxy/internal/clients/ratelimiter_service"

	"github.com/google/uuid"
)

type ProxyHandler struct {
	configs           map[string]*Resource
	transport         http.RoundTripper
	rateLimiterClient *ratelimiterservice.RateLimiterClient
	cacherClient      *cacherservice.CacherClient
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
		cacherClient:      cacherservice.NewCacherClient(cfg.CacherURL),
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

	cacheKey, _ := ph.cacherClient.GenerateCacheKey(r)
	cachedData, err := ph.cacherClient.GetCache(ctx, cacheKey)
	if err == nil {
		log.Printf("request with id=%s was cached", requestID)

		w.Write([]byte(cachedData))
		return
	} else {
		log.Printf("error while getting cache by key=%s: %v", cacheKey, err)
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
	respBody, _ := io.ReadAll(resp.Body)

	// TODO: раскомментить и пофиксить когда-нибудь
	// go func() {
		err = ph.cacherClient.SetCache(ctx, cacheKey, string(respBody))
		if err != nil {
			log.Printf("failed to cache response: %v", err)
		}
	// }()

	for k, v := range resp.Header {
		w.Header()[k] = v
	}

	w.WriteHeader(resp.StatusCode)

	w.Write(respBody)
}

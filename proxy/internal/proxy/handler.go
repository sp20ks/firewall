package proxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	cacher "proxy/internal/clients/cacher_service"
	ratelimiter "proxy/internal/clients/ratelimiter_service"
	rules "proxy/internal/clients/rules_engine_service"
	"proxy/internal/config"

	"github.com/google/uuid"
)

type ResourceMap map[string]map[string]rules.Resource

type ProxyHandler struct {
	resources         ResourceMap
	transport         http.RoundTripper
	rateLimiterClient *ratelimiter.RateLimiterClient
	cacherClient      *cacher.CacherClient
	rulesEngineClient *rules.RulesEngineClient
}

func NewProxyHandler(cfg *config.Config) (*ProxyHandler, error) {
	rulesClient := rules.NewRulesEngineClient(cfg.RulesEngineURL)

	resources, err := rulesClient.GetResources()
	if err != nil {
		return nil, fmt.Errorf("failed to load resources from rules engine: %w", err)
	}

	// TODO: доступные ресурсы должны обновляться постоянно, а не только при инициализации хэндлера, т.е. при деплое
	resourcesMap := make(ResourceMap)
	for _, res := range resources {
		if _, exists := resourcesMap[res.URL]; !exists {
			resourcesMap[res.URL] = make(map[string]rules.Resource)
		}
		resourcesMap[res.URL][res.Method] = res
	}

	return &ProxyHandler{
		resources: resourcesMap,
		transport: &http.Transport{
			DisableKeepAlives: true,
		},
		rateLimiterClient: ratelimiter.NewRateLimiterClient(cfg.RateLimiterURL),
		cacherClient:      cacher.NewCacherClient(cfg.CacherURL),
		rulesEngineClient: rulesClient,
	}, nil
}

func (ph *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := uuid.NewString()
	ctx = context.WithValue(ctx, "request-id", requestID)

	resourceMethods, pathExists := ph.resources[r.URL.Path]
	if !pathExists {
		errResp := ErrorResponse{
			Error:      "endpoint not found",
			StatusCode: http.StatusNotFound,
			Timestamp:  time.Now().Format(time.RFC3339),
			RequestID:  requestID,
		}
		WriteJSONResponse(w, errResp, http.StatusNotFound)
		return
	}

	resource, methodExists := resourceMethods[r.Method]
	if !methodExists {
		errResp := ErrorResponse{
			Error:      "method not allowed",
			StatusCode: http.StatusMethodNotAllowed,
			Timestamp:  time.Now().Format(time.RFC3339),
			RequestID:  requestID,
		}
		WriteJSONResponse(w, errResp, http.StatusMethodNotAllowed)
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

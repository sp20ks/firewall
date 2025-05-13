package proxy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"proxy/internal/config"
	"proxy/internal/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
	cacher "proxy/internal/clients/cacher_service"
	ratelimiter "proxy/internal/clients/ratelimiter_service"
	rules "proxy/internal/clients/rules_engine_service"
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
	l := logger.Logger()

	resourceMethods, pathExists := ph.resources[r.URL.Path]
	if !pathExists {
		WriteJSONResponse(w, NewErrorResponse("endpoint not found", http.StatusNotFound, requestID), http.StatusNotFound)
		return
	}

	resource, methodExists := resourceMethods[r.Method]
	if !methodExists {
		WriteJSONResponse(w, NewErrorResponse("method not allowed", http.StatusMethodNotAllowed, requestID), http.StatusMethodNotAllowed)
		return
	}

	if code, err := ph.validateRequest(r); err != nil {
		WriteJSONResponse(w, NewErrorResponse(err.Error(), code, requestID), code)
		return
	}
	cacheKey, _ := ph.cacherClient.GenerateCacheKey(r)
	cachedData, err := ph.cacherClient.GetCache(ctx, cacheKey)
	if err == nil {
		l.Info("request was cached", zap.String("request_id", requestID))

		w.Write([]byte(cachedData))
		return
	} else {
		l.Info("error while getting cache by key", zap.String("key", cacheKey), zap.Error(err))
	}

	req, err := ph.modifyRequest(ctx, r, resource)
	if err != nil {
		WriteJSONResponse(w, NewErrorResponse("internal server error", http.StatusInternalServerError, requestID), http.StatusInternalServerError)
		return
	}

	resp, err := ph.forwardRequest(ctx, req)
	if err != nil {
		l.Info("error while proxing request", zap.String("key", cacheKey), zap.Error(err))

		WriteJSONResponse(w, NewErrorResponse("proxy error", http.StatusInternalServerError, requestID), http.StatusInternalServerError)
		return
	}
	respBody, _ := io.ReadAll(resp.Body)

	// TODO: раскомментить и пофиксить когда-нибудь
	// go func() {
	err = ph.cacherClient.SetCache(ctx, cacheKey, string(respBody))
	if err != nil {
		l.Info("failed to cache response", zap.String("key", cacheKey), zap.Error(err))
	}
	// }()

	for k, v := range resp.Header {
		w.Header()[k] = v
	}

	w.WriteHeader(resp.StatusCode)

	w.Write(respBody)
}

package proxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.uber.org/zap"
	rules "proxy/internal/clients/rules_engine_service"
	"proxy/internal/logger"
)

func (ph *ProxyHandler) modifyRequest(ctx context.Context, r *http.Request, resource rules.Resource) (*http.Request, error) {
	rawUrl, err := url.Parse(resource.Host)
	if err != nil {
		return nil, err
	}

	queryParams := r.URL.Query().Encode()
	url := fmt.Sprintf("%s?%s", rawUrl, queryParams)

	req, err := http.NewRequestWithContext(ctx, r.Method, url, r.Body)
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

func (ph *ProxyHandler) validateRequest(r *http.Request) (int, error) {
	ip := ReadUserIP(r)
	l := logger.Logger()

	if allowed, err := ph.rateLimiterClient.CheckLimit(ip); !allowed {
		l.Info("rate limit error", zap.String("ip", ip), zap.Error(err))
		return http.StatusTooManyRequests, fmt.Errorf("too many requests")
	}

	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = io.ReadAll(r.Body)
	}

	headers := make(map[string]string)
	for k, v := range r.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	analysisResp, err := ph.rulesEngineClient.AnalyzeRequest(
		ip,
		r.Method,
		r.URL.String(),
		string(bodyBytes),
		headers,
	)
	if err != nil {
		l.Info("error analyzing request", zap.Error(err))
		return http.StatusInternalServerError, fmt.Errorf("error analyzing request")
	}

	switch analysisResp.Action {
	case "block":
		l.Info("blocker request from ip", zap.String("ip", ip), zap.String("reason", analysisResp.Reason))
		return http.StatusForbidden, fmt.Errorf("request blocked: %s", analysisResp.Reason)
	case "allow":
		if analysisResp.ModifiedBody != "" {
			r.Body = io.NopCloser(bytes.NewBufferString(analysisResp.ModifiedBody))
		}
		if analysisResp.ModifiedURL != "" {
			modifiedURL, err := url.Parse(analysisResp.ModifiedURL)
			if err == nil {
				r.URL = modifiedURL
			}
		}
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

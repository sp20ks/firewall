package proxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type ProxyConfig struct {
	Resources []Resource
}

type Resource struct {
	Host     string
	Endpoint string
}

type ProxyHandler struct {
	configs   map[string]*Resource
	transport http.RoundTripper
}

func NewProxyHandler(resources []Resource) (*ProxyHandler, error) {
	configs := make(map[string]*Resource)

	for _, r := range resources {
		url, err := url.Parse(r.Host)
		if err != nil {
			return nil, fmt.Errorf("invalid URL %s: %v", r.Host, err)
		}

		configs[r.Endpoint] = &Resource{
			Host:     url.String(),
			Endpoint: r.Endpoint,
		}
	}

	return &ProxyHandler{
		configs: configs,
		transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}, nil
}

func (ph *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resource, ok := ph.configs[r.URL.Path]
	if !ok {
		http.Error(w, "endpoint not found", http.StatusNotFound)
		return
	}

	if err := validateRequest(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req, err := ph.modifyRequest(ctx, r, resource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := ph.forwardRequest(ctx, req)
	if err != nil {
		http.Error(w, "proxy error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	copyResponse(w, resp)
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

// unused func
func validateRequest(r *http.Request) error {
	return nil
}

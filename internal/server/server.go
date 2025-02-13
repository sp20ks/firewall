package server

import (
	"firewall/internal/config"
	"firewall/internal/proxy"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func Start(cfg *config.Config) error {
	addr := cfg.HTTPServer.Address

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)

	for _, resource := range cfg.Resources {
		url, _ := url.Parse(resource.Host)
		mux.HandleFunc(resource.Endpoint, proxy.ProxyRequestHandler(url, resource.Endpoint))
	}

	log.Printf("Starting proxy server on %s\n", addr)

	if err := http.ListenAndServe(cfg.HTTPServer.Address, mux); err != nil {
		return fmt.Errorf("error starting server: %v", err)
	}
	return nil
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

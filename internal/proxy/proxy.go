package proxy

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var customTransport = http.DefaultTransport

func ProxyRequestHandler(url *url.URL, endpoint string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Receive request: %v", r)

		r.URL.Host = url.Host
		r.URL.Scheme = url.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = url.Host
		path := r.URL.Path
		r.URL.Path = strings.TrimLeft(path, endpoint)

		log.Printf("Proxing request: %v", r)
		resp, err := customTransport.RoundTrip(r)
		if err != nil {
			log.Printf("error sending proxy request: %v", err)
			http.Error(w, "error sending proxy request", http.StatusInternalServerError)
			return
		}
		log.Printf("Received response: %v", resp)
		defer resp.Body.Close()

		for name, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		w.WriteHeader(resp.StatusCode)

		io.Copy(w, resp.Body)
	}
}

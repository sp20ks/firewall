package proxy

import (
	"fmt"
	"net/url"
)

type ProxyConfig struct {
	Resources []Resource
}

type Resource struct {
	Host     string
	Endpoint string
}

func loadConfigs(resources []Resource) (map[string]*Resource, error) {
	configs := make(map[string]*Resource)

	for _, r := range resources {
		parsedURL, err := url.Parse(r.Host)
		if err != nil {
			return nil, fmt.Errorf("invalid URL %s: %v", r.Host, err)
		}

		configs[r.Endpoint] = &Resource{
			Host:     parsedURL.String(),
			Endpoint: r.Endpoint,
		}
	}

	return configs, nil
}

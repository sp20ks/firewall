package rulesengineservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Resource struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Host   string `json:"host"`
	Method string `json:"http_method"`
}

type ResourcesResponse struct {
	Resources []Resource `json:"resources"`
}

type RulesEngineClient struct {
	rulesEngineURL string
}

func NewRulesEngineClient(url string) *RulesEngineClient {
	return &RulesEngineClient{rulesEngineURL: url}
}

func (re *RulesEngineClient) GetResources() ([]Resource, error) {
	request, err := http.NewRequest(http.MethodGet, re.rulesEngineURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{Timeout: 2 * time.Second}

	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error requesting resource list: %w", err)
	}
	defer resp.Body.Close()

	var resourcesResp ResourcesResponse
	if err := json.NewDecoder(resp.Body).Decode(&resourcesResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return resourcesResp.Resources, nil
}

package rulesengineservice

import (
	"bytes"
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

type AnalizerRequest struct {
	IP      string            `json:"ip"`
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type AnalizerResponse struct {
	Action       string `json:"action"`
	ModifiedURL  string `json:"modified_url,omitempty"`
	ModifiedBody string `json:"modified_body,omitempty"`
	Reason       string `json:"reason"`
}

func NewRulesEngineClient(url string) *RulesEngineClient {
	return &RulesEngineClient{rulesEngineURL: url}
}

func (re *RulesEngineClient) GetResources() ([]Resource, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/resources", re.rulesEngineURL), nil)
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

func (re *RulesEngineClient) AnalizeRequest(ip, method, url, body string, headers map[string]string) (*AnalizerResponse, error) {
	respBody, err := json.Marshal(AnalizerRequest{IP: ip, Method: method, URL: url, Body: body, Headers: headers})
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/analize", re.rulesEngineURL), bytes.NewBuffer(respBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error check request: %w", err)
	}
	defer resp.Body.Close()
	var analizerResp AnalizerResponse
	if err := json.NewDecoder(resp.Body).Decode(&analizerResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &analizerResp, nil
}

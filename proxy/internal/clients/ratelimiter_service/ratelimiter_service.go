package ratelimiterservice

import (
	"fmt"
	"net/http"
	"time"
)

type RateLimiterClient struct {
	rateLimiterURL string
}

func NewRateLimiterClient(url string) *RateLimiterClient {
	return &RateLimiterClient{rateLimiterURL: url}
}

func (rl *RateLimiterClient) CheckLimit(ip string) (bool, error) {
	url := fmt.Sprintf("%s?ip=%s", rl.rateLimiterURL, ip)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{Timeout: 2 * time.Second}

	resp, err := client.Do(request)
	if err != nil {
		return false, fmt.Errorf("error requesting rate limiter: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return false, nil
	} else if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected response from rate limiter: %d", resp.StatusCode)
	}

	return true, nil
}

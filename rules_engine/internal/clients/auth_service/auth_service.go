package authservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type AuthClient struct {
	authURL string
}

func NewAuthClient(url string) *AuthClient {
	return &AuthClient{authURL: url}
}

func (a *AuthClient) VerifyToken(token string) (*UserResponse, error) {
	url := fmt.Sprintf("%s?token=%s", a.authURL, token)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{Timeout: 2 * time.Second}

	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error verify token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("token is invalid")
	}

	if resp.StatusCode == http.StatusBadRequest {
		return nil, fmt.Errorf("request is invalid")
	}

	var userResp UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &userResp, nil
}

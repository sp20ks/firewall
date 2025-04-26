package cacherservice

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CacherClient struct {
	cacherURL string
	client    *http.Client
}

type CacherRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewCacherClient(url string) *CacherClient {
	return &CacherClient{
		cacherURL: url,
		client:    &http.Client{Timeout: 2 * time.Second},
	}
}

func (cc *CacherClient) GenerateCacheKey(req *http.Request) (string, error) {
	encodedQuery := req.URL.Query().Encode()
	key := fmt.Sprintf("%s:%s?%s", req.Method, req.URL.Path, encodedQuery)

	if req.Method == http.MethodPost || req.Method == http.MethodPut || req.Method == http.MethodPatch {
		bodyHash, err := hashRequestBody(req)
		if err != nil {
			return "", err
		}
		key += fmt.Sprintf(":%s", bodyHash)
	}

	return key, nil
}

func (cc *CacherClient) GetCache(ctx context.Context, key string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s?key=%s", cc.cacherURL, key), nil)
	if err != nil {
		return "", err
	}
	resp, err := cc.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return "", fmt.Errorf("cache not found")
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result CacherRequest
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Value, nil
}

func (cc *CacherClient) SetCache(ctx context.Context, key, value string) error {
	body, err := json.Marshal(CacherRequest{Key: key, Value: value})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cc.cacherURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := cc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to set cache: %d", resp.StatusCode)
	}

	return nil
}

func hashRequestBody(req *http.Request) (string, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	hash := sha256.Sum256(body)
	return hex.EncodeToString(hash[:]), nil
}

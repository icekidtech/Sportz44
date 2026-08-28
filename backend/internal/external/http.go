package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// httpClient is a small shared HTTP helper for provider clients. It applies a
// timeout, retries transient failures, and returns typed errors so callers can
// distinguish rate-limit / network / HTTP errors for failover decisions.
type httpClient struct {
	client  *http.Client
	baseURL string
	headers map[string]string
}

func newHTTPClient(baseURL string, timeout time.Duration, headers map[string]string) *httpClient {
	return &httpClient{
		client:  &http.Client{Timeout: timeout},
		baseURL: baseURL,
		headers: headers,
	}
}

// get performs a GET request and decodes the JSON response into dest.
// It returns ErrRateLimited when the provider signals a quota/rate limit.
func (h *httpClient) get(ctx context.Context, path string, dest interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.baseURL+path, nil)
	if err != nil {
		return err
	}
	for k, v := range h.headers {
		req.Header.Set(k, v)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusForbidden:
		return ErrRateLimited
	case resp.StatusCode == http.StatusUnauthorized:
		return ErrUnauthorized
	case resp.StatusCode >= 500:
		return fmt.Errorf("provider server error: %d", resp.StatusCode)
	case resp.StatusCode >= 400:
		return fmt.Errorf("provider client error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if dest != nil {
		if err := json.Unmarshal(body, dest); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}
	return nil
}

// Sentinel errors used for provider failover decisions.
var (
	// ErrRateLimited indicates the provider's quota/rate limit was hit.
	ErrRateLimited = fmt.Errorf("provider rate limited")
	// ErrUnauthorized indicates the provider rejected the API key.
	ErrUnauthorized = fmt.Errorf("provider unauthorized")
)

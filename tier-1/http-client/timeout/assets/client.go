package assets

import (
	"context"
	"net/http"
	"time"
)

type HttpClient interface {
	Get(url string) (*http.Response, error)
}

type SecureHttpClient struct {
	client  *http.Client
	timeout time.Duration
}

func NewSecureHttpClient(timeout time.Duration) *SecureHttpClient {
	return &SecureHttpClient{
		client:  &http.Client{},
		timeout: timeout,
	}
}

func (c *SecureHttpClient) Get(ctx context.Context, url string) (*http.Response, error) {
	// Create a new context with a timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel() // Ensure the cancel function is called to release resources

	req, err := http.NewRequestWithContext(ctxWithTimeout, http.MethodGet, url, nil)
	if err != nil {
		cancel()
		return nil, err
	}

	return c.client.Do(req)
}

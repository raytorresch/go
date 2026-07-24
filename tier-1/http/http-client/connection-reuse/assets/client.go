package assets

import (
	"context"
	"io"
	"net/http"
	"net/http/httptrace"
	"sync/atomic"
	"time"
)

type HttpClient interface {
	Get(ctx context.Context, url string) (*http.Response, error)
}

type SecureHttpClient struct {
	client  *http.Client
	timeout time.Duration
	Metrics *MetricConnection
}

type MetricConnection struct {
	NewCon   int64
	ReuseCon int64
}

func (m *MetricConnection) Reset() {
	atomic.StoreInt64(&m.NewCon, 0)
	atomic.StoreInt64(&m.ReuseCon, 0)
}

func NewTransport(idleCont, indleHostCon int, idleTimeout time.Duration) *http.Transport {
	return &http.Transport{
		MaxIdleConns:        idleCont,
		MaxIdleConnsPerHost: indleHostCon,
		IdleConnTimeout:     idleTimeout,
	}
}
func NewSecureHttpClient(timeout time.Duration, transport *http.Transport, metrics *MetricConnection) *SecureHttpClient {
	return &SecureHttpClient{
		client: &http.Client{
			Transport: transport,
			Timeout:   timeout,
		},
		timeout: timeout, //used to manage manually time out, close wrapper is necesary doing this
		Metrics: metrics,
	}
}

func (c *SecureHttpClient) Get(ctx context.Context, url string) (*http.Response, error) {
	// Create a new context with a timeout
	// ctxWithTimeout, cancel := context.WithTimeout(ctx, c.timeout)

	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			if connInfo.Reused {
				atomic.AddInt64(&c.Metrics.ReuseCon, 1)
			} else {
				atomic.AddInt64(&c.Metrics.NewCon, 1)
			}
		},
	}

	// ctxWithTrace := httptrace.WithClientTrace(ctxWithTimeout, trace)
	ctxWithTrace := httptrace.WithClientTrace(ctx, trace)

	req, err := http.NewRequestWithContext(ctxWithTrace, http.MethodGet, url, nil)
	if err != nil {
		// cancel()
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		// cancel()
		return nil, err
	}

	// Don't cancel until the body is closed: canceling right after Do()
	// returns (while the caller still has to read/close the body) would
	// force the Transport to tear down the connection instead of
	// returning it to the idle pool, killing reuse entirely.
	// res.Body = &cancelOnCloseBody{ReadCloser: res.Body, cancel: cancel}
	return res, nil
}

type cancelOnCloseBody struct {
	io.ReadCloser
	cancel context.CancelFunc
}

// Close method wraper
// Close the body and then cancel the context
func (b *cancelOnCloseBody) Close() error {
	err := b.ReadCloser.Close()
	b.cancel()
	return err
}

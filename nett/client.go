package nett

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/go-dj/dj"
)

// Client is a REST client for the given base URL.
type Client struct {
	base string
	rt   http.RoundTripper

	preHooks     []func(*http.Request)
	preHooksLock sync.RWMutex

	postHooks     []func(*http.Response)
	postHooksLock sync.RWMutex
}

// NewClient returns a new client for the given base URL.
// It uses the provided round tripper to send HTTP requests.
func NewClient(baseURL string, rt http.RoundTripper) *Client {
	return &Client{
		base: baseURL,
		rt:   rt,
	}
}

// AddPreReqHook adds a hook that is called before the request is sent.
func (c *Client) AddPreReqHook(hook func(*http.Request)) {
	c.preHooksLock.Lock()
	defer c.preHooksLock.Unlock()

	c.preHooks = append(c.preHooks, hook)
}

// AddPostRespHook adds a hook that is called after the response is received.
func (c *Client) AddPostReqHook(hook func(*http.Response)) {
	c.postHooksLock.Lock()
	defer c.postHooksLock.Unlock()

	c.postHooks = append(c.postHooks, hook)
}

// Get makes a GET request to the given path.
func Get[Res any](ctx context.Context, c *Client, path string) (Res, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.base+path, nil)
	if err != nil {
		return dj.Zero[Res](), fmt.Errorf("failed to create request: %w", err)
	}

	return Do[Res](c, req)
}

// Post makes a POST request to the given path.
func Post[Req, Res any](ctx context.Context, c *Client, path string, req Req) (Res, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return dj.Zero[Res](), fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+path, bytes.NewReader(b))
	if err != nil {
		return dj.Zero[Res](), fmt.Errorf("failed to create request: %w", err)
	}

	return Do[Res](c, httpReq)
}

// Put makes a PUT request to the given path.
func Put[Req, Res any](ctx context.Context, c *Client, path string, req Req) (Res, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return dj.Zero[Res](), fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, c.base+path, bytes.NewReader(body))
	if err != nil {
		return dj.Zero[Res](), fmt.Errorf("failed to create request: %w", err)
	}

	return Do[Res](c, httpReq)
}

// Delete makes a DELETE request to the given path.
func Delete[Res any](ctx context.Context, c *Client, path string) (Res, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.base+path, nil)
	if err != nil {
		return dj.Zero[Res](), fmt.Errorf("failed to create request: %w", err)
	}

	return Do[Res](c, req)
}

// Do sends the given request and returns the response body and an error, if any.
func Do[Res any](c *Client, req *http.Request) (Res, error) {
	c.preHooksLock.RLock()
	defer c.preHooksLock.RUnlock()

	for _, hook := range c.preHooks {
		hook(req)
	}

	resp, err := c.rt.RoundTrip(req)
	if err != nil {
		return dj.Zero[Res](), fmt.Errorf("failed to send request: %w", err)
	}

	c.postHooksLock.RLock()
	defer c.postHooksLock.RUnlock()

	for _, hook := range c.postHooks {
		hook(resp)
	}

	if resp.StatusCode != http.StatusOK {
		return dj.Zero[Res](), fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return dj.Zero[Res](), fmt.Errorf("failed to read response body: %w", err)
	}

	var res Res

	if err := json.Unmarshal(b, &res); err != nil {
		return dj.Zero[Res](), fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return res, nil
}

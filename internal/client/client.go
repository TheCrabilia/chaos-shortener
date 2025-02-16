package client

import (
	"bytes"
	"fmt"
	"net/http"
)

type Client struct {
	*http.Client
	baseURL string
}

type Serializable interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

func New(baseURL string) *Client {
	return &Client{
		Client:  &http.Client{},
		baseURL: baseURL,
	}
}

func (c *Client) Request(method, endpoint string, body Serializable) (*http.Response, error) {
	b, err := body.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	req, err := http.NewRequest(method, c.baseURL+endpoint, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.Do(req)
}

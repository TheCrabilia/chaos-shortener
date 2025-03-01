package client

import (
	"fmt"
	"io"
	"net/http"

	"github.com/TheCrabilia/chaos-shortener/internal/server/api"
)

type ShortenURLOpts struct {
	URL    string
	Repeat int
}

func (c *Client) ShortenURL(opts *ShortenURLOpts) ([]string, error) {
	var (
		results []string
		errCh   = make(chan error)
	)

	req := &api.ShortenRequest{
		URL: opts.URL,
	}

	for range opts.Repeat {
		resp, err := c.Request(http.MethodPost, "/shorten", req)
		if err != nil {
			errCh <- fmt.Errorf("failed to send request: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			errCh <- fmt.Errorf("failed to read response body: %w", err)
		}

		var data api.ShortenResponse
		if err := data.Unmarshal(body); err != nil {
			errCh <- fmt.Errorf("failed to unmarshal response: %w", err)
		}

		results = append(results, data.ShortURL)
	}

	return results, nil
}

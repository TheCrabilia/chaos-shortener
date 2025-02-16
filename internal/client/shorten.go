package client

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/TheCrabilia/chaos-shortener/internal/server/api"
)

type ShortenURLOpts struct {
	URL      string
	Repeat   int
	Parallel bool
}

func (c *Client) ShortenURL(opts *ShortenURLOpts) ([]string, error) {
	var (
		results []string
		errCh   = make(chan error)
		mu      sync.Mutex
		wg      sync.WaitGroup
	)

	req := &api.ShortenRequest{
		URL: opts.URL,
	}

	reqFunc := func() {
		defer wg.Done()
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

		mu.Lock()
		defer mu.Unlock()

		results = append(results, data.ShortURL)
	}

	for i := 0; i < opts.Repeat; i++ {
		wg.Add(1)

		if opts.Parallel {
			go reqFunc()
		} else {
			reqFunc()
		}
	}

	wg.Wait()

	return results, nil
}

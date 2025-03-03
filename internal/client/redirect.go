package client

import "net/http"

type RedirectURLOpts struct {
	ID string
}

func (c *Client) RedirectURL(opts *RedirectURLOpts) error {
	_, err := c.Request(http.MethodGet, "/r/"+opts.ID, nil)

	return err
}

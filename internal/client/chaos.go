package client

import (
	"net/http"

	"github.com/TheCrabilia/chaos-shortener/internal/server/api"
)

type UpdateChaosSettingsOpts struct {
	LatencyRate float64
	ErrorRate   float64
	OutageRate  float64
}

func (c *Client) UpdateChaosSettings(opts *UpdateChaosSettingsOpts) {
	req := &api.InjectorRequest{
		LatencyRate: opts.LatencyRate,
		ErrorRate:   opts.ErrorRate,
		OutageRate:  opts.OutageRate,
	}

	if _, err := c.Request(http.MethodPost, "/chaos", req); err != nil {
		panic(err)
	}
}

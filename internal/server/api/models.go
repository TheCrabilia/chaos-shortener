package api

import "encoding/json"

type ShortenRequest struct {
	URL string `json:"url"`
}

func (sr *ShortenRequest) Marshal() ([]byte, error) {
	b, err := json.Marshal(sr)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (sr *ShortenRequest) Unmarshal(b []byte) error {
	return json.Unmarshal(b, sr)
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

func (sr *ShortenResponse) Marshal() ([]byte, error) {
	b, err := json.Marshal(sr)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (sr *ShortenResponse) Unmarshal(b []byte) error {
	return json.Unmarshal(b, sr)
}

type InjectorRequest struct {
	LatencyRate  float64 `json:"latency_rate"`
	ErrorRate    float64 `json:"error_rate"`
	ConnDropRate float64 `json:"conn_drop_rate"`
	OutageRate   float64 `json:"outage_rate"`
}

func (ir *InjectorRequest) Marshal() ([]byte, error) {
	b, err := json.Marshal(ir)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (ir *InjectorRequest) Unmarshal(b []byte) error {
	return json.Unmarshal(b, ir)
}

type InjectorResponse struct {
	LatencyRate  float64 `json:"latency_rate"`
	ErrorRate    float64 `json:"error_rate"`
	ConnDropRate float64 `json:"conn_drop_rate"`
	OutageRate   float64 `json:"outage_rate"`
}

func (ir *InjectorResponse) Marshal() ([]byte, error) {
	b, err := json.Marshal(ir)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (ir *InjectorResponse) Unmarshal(b []byte) error {
	return json.Unmarshal(b, ir)
}

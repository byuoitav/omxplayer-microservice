package data

import (
	"context"
)

type ConfigService interface {
	GetStreamConfig(ctx context.Context, streamURL string) (Stream, error)
}

type StreamConfig struct {
	Streams map[string]Stream `json:"streams"`
}

type Stream struct {
	URL         string `json:"url"`
	Secret      string `json:"secret"`
	QueryPrefix string `json:"queryPrefix"`
	Duration    string `json:"duration"`
}

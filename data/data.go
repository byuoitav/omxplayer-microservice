package data

import (
	"context"
)

type ConfigService interface {
	GetStreamConfig(ctx context.Context, streamURL string) (Stream, error)
}

type Stream struct {
	Secret      string `json:"secret"`
	QueryPrefix string `json:"queryPrefix"`
	Duration    string `json:"duration"`
}

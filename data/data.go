package data

import (
	"context"
)

type ConfigService interface {
	GetStreamConfig(ctx context.Context, streamURL string) (Stream, error)
	GetDeviceConfig(ctx context.Context, hostanem string) (Device, error)
}

type Stream struct {
	URL         string `json:"url"`
	Secret      string `json:"secret"`
	QueryPrefix string `json:"queryPrefix"`
	Duration    string `json:"duration"`
}

type Device struct {
	Args []string `json:"args"`
}

package cache

import (
	"context"

	"github.com/byuoitav/omxplayer-microservice/data"
)

type ConfigService struct {
	ConfigService data.ConfigService
}

func (c *ConfigService) GetStreamConfig(ctx context.Context, stream string) (data.Stream, error) {
	// TODO cache
	return c.ConfigService.GetStreamConfig(ctx, stream)
}

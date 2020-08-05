package cache

import (
	"context"

	"github.com/byuoitav/omxplayer-microservice/data"
)

type ConfigService struct {
	ConfigService data.ConfigService
}

func (c *ConfigService) GetStreamConfig(ctx context.Context, streamURL string) (data.Stream, error) {
	stream, err := c.ConfigService.GetStreamConfig(ctx, streamURL)
	if err != nil {
		//check the cache
	}

	//store in cache
	return stream, nil
}

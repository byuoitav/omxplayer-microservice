package couch

import (
	"context"
	"fmt"

	"github.com/byuoitav/omxplayer-microservice/data"
	_ "github.com/go-kivik/couchdb/v3"
	kivik "github.com/go-kivik/kivik/v3"
)

const (
	_streamConfigDoc = "streams"
)

type ConfigService struct {
	Client         *kivik.Client
	StreamConfigDB string
}

type streamConfig struct {
	Streams map[string]data.Stream `json:"streams"`
}

func (c *ConfigService) GetStreamConfig(ctx context.Context, streamURL string) (data.Stream, error) {
	var config streamConfig

	db := c.Client.DB(ctx, c.StreamConfigDB)
	if err := db.Get(ctx, "streams").ScanDoc(&config); err != nil {
		return data.Stream{}, fmt.Errorf("unable to get stream configs: %w", err)
	}

	return config.Streams[streamURL], nil
}

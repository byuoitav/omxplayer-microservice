package couch

import (
	"context"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/omxplayer-microservice/data"
	_ "github.com/go-kivik/couchdb/v3"
	kivik "github.com/go-kivik/kivik/v3"
)

type ConfigService struct {
	Client         *kivik.Client
	StreamConfigDB string
}

func (c *ConfigService) GetStreamConfig(ctx context.Context, streamURL string) (data.Stream, error) {
	var config data.StreamConfig

	db := c.Client.DB(ctx, c.StreamConfigDB)
	if err := db.Get(ctx, "streams").ScanDoc(&config); err != nil {
		log.L.Errorf("error getting stream config doc: %s", err)
		return data.Stream{}, err
	}

	if s, ok := config.Streams[streamURL]; ok {
		return s, nil
	}

	return data.Stream{}, nil
}

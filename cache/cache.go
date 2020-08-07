package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/omxplayer-microservice/data"
	bolt "go.etcd.io/bbolt"
)

const (
	_bucket = "streams"
)

type configService struct {
	configService data.ConfigService
	db            *bolt.DB
}

func New(cs data.ConfigService, path string) (data.ConfigService, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to open cache: %w", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(_bucket))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to initialize cache: %w", err)
	}

	return &configService{
		configService: cs,
		db:            db,
	}, nil
}

func (c *configService) GetStreamConfig(ctx context.Context, streamURL string) (data.Stream, error) {
	stream, err := c.configService.GetStreamConfig(ctx, streamURL)
	if err != nil {
		return c.streamConfigFromCache(ctx, streamURL)
	}

	if err := c.cacheStream(ctx, streamURL, stream); err != nil {
		log.L.Warnf("unable to cache stream %q: %s", streamURL, err)
	}

	return stream, nil
}

func (c *configService) cacheStream(ctx context.Context, streamURL string, stream data.Stream) error {
	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_bucket))
		if b == nil {
			return fmt.Errorf("stream bucket does not exist")
		}

		bytes, err := json.Marshal(stream)
		if err != nil {
			return fmt.Errorf("unable to marshal stream: %w", err)
		}

		if err = b.Put([]byte(streamURL), bytes); err != nil {
			return fmt.Errorf("unable to put stream: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *configService) streamConfigFromCache(ctx context.Context, streamURL string) (data.Stream, error) {
	var stream data.Stream

	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_bucket))
		if b == nil {
			return fmt.Errorf("stream bucket does not exist")
		}

		bytes := b.Get([]byte(streamURL))
		if bytes == nil {
			return fmt.Errorf("stream not in cache")
		}

		if err := json.Unmarshal(bytes, &stream); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return stream, fmt.Errorf("unable to get stream from cache: %w", err)
	}

	return stream, nil
}

package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/omxplayer-microservice/data"
	bolt "go.etcd.io/bbolt"
)

type StreamCache struct {
	Streams []data.Stream
}

type ConfigService struct {
	ConfigService data.ConfigService
	DB            *bolt.DB
}

func (c *ConfigService) GetStreamConfig(ctx context.Context, streamURL string) (data.Stream, error) {
	stream, err := c.ConfigService.GetStreamConfig(ctx, streamURL)
	if err != nil {
		//check the cache
		err := c.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("STREAMS"))
			if b == nil {
				return fmt.Errorf("stream bucket does not exist")
			}

			bytes := b.Get([]byte(streamURL))
			if bytes == nil {
				return fmt.Errorf("employee not in cache")
			}

			if err := json.Unmarshal(bytes, &stream); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			log.L.Errorf("unable to find stream in cache: %s", err.Error())
			return data.Stream{}, fmt.Errorf("unable to find stream in cache: %s", err.Error())
		}

		return stream, nil
	}

	//store in cache
	err = c.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("STREAMS"))

		bytes, err := json.Marshal(stream)
		if err != nil {
			log.L.Errorf("unable to marshal stream %s: %s", streamURL, err.Error())
			return fmt.Errorf("unable to marshal stream %s: %s", streamURL, err.Error())
		}

		if err = b.Put([]byte(streamURL), bytes); err != nil {
			log.L.Errorf("unable to cache stream %s: %s", streamURL, err.Error())
			return fmt.Errorf("unable to cache stream %s: %s", streamURL, err.Error())
		}

		return nil
	})
	if err != nil {
		return data.Stream{}, err
	}

	return stream, nil
}

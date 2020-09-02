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
	_streamBucket = "streams"
	_deviceBucket = "devices"
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
		_, err := tx.CreateBucketIfNotExists([]byte(_streamBucket))
		if err != nil {
			return err
		}

		_, err := tx.CreateBucketIfNotExists([]byte(_deviceBucket))
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
		stream, cacheErr := c.streamConfigFromCache(ctx, streamURL)
		if cacheErr != nil {
			log.L.Warnf("unable to get stream %q from cache: %s", streamURL, cacheErr)
			return stream, err
		}

		return stream, nil
	}

	if err := c.cacheStream(ctx, streamURL, stream); err != nil {
		log.L.Warnf("unable to cache stream %q: %s", streamURL, err)
	}

	return stream, nil
}

func (c *configService) cacheStream(ctx context.Context, streamURL string, stream data.Stream) error {
	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_streamBucket))
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
		b := tx.Bucket([]byte(_streamBucket))
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
		return stream, err
	}

	return stream, nil
}

func (c *configService) GetDeviceConfig(ctx context.Context, hostname string) (data.Device, error) {
	device, err := c.configService.GetDeviceConfig(ctx, hostname)
	if err != nil {
		device, cacheErr := c.deviceConfigFromCache(ctx, hostname)
		if cacheErr != nil {
			log.L.Warnf("unable to get device %s from cache: %s", hostname, cacheErr)
			return device, err
		}

		return device, nil
	}

	if err := c.cacheDevice(ctx, hostname, device); err != nil {
		log.L.Warnf("unable to cache device %s: %s", hostname, err)
	}

	return device, nil
}

func (c *configService) cacheDevice(ctx context.Context, hostname string, device data.Device) error {
	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_deviceBucket))
		if b == nil {
			return fmt.Errorf("device bucket does not exist")
		}

		bytes, err := json.Marshal(device)
		if err != nil {
			return fmt.Errorf("unable to marshal device")
		}

		if err = b.Put([]byte(hostname), bytes); err != nil {
			return fmt.Errorf("unable to put device: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *configService) deviceConfigFromCache(ctx context.Context, hostname) (data.Device, error) {
	var device data.Device

	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_deviceBucket))
		if b == nil {
			return fmt.Errorf("device bucket does not exist")
		}

		bytes := b.Get([]byte(hostname))
		if bytes == nil {
			return fmt.Errorf("device not in cache")
		}

		if err := json.Unmarshal(bytes, &device); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return device, err
	}

	return device, nil
}
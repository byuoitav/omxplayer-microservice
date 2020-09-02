package cache

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/byuoitav/omxplayer-microservice/data"
	"github.com/stretchr/testify/require"
)

type mockConfigService struct {
	streams map[string]data.Stream
	devices map[string]data.Device
}

func (m *mockConfigService) GetStreamConfig(ctx context.Context, streamURL string) (data.Stream, error) {
	stream, ok := m.streams[streamURL]
	if !ok {
		return stream, fmt.Errorf("stream not found")
	}

	return stream, nil
}

func (m *mockConfigService) GetDeviceConfig(ctx context.Context, hostname string) (data.Device, error) {
	device, ok := m.devices[hostname]
	if !ok {
		return device, fmt.Errorf("device not found")
	}

	return device, nil
}

func TestCache(t *testing.T) {
	hiStream := data.Stream{
		Secret:      "hi",
		QueryPrefix: "hello",
		Duration:    "1s",
	}
	hiDevice := data.Device{
		Args: []string{"hi", "device"},
	}
	mock := &mockConfigService{
		streams: map[string]data.Stream{
			"hi.com": hiStream,
		},
		devices: map[string]data.Device{
			"hiDevice": hiDevice,
		},
	}

	file := os.TempDir() + "/omxplayer-microservice-cache-test.db"
	config, err := New(mock, file)
	require.NoError(t, err)
	defer os.Remove(file)

	t.Run("StreamPassthrough", func(t *testing.T) {
		stream, err := config.GetStreamConfig(context.TODO(), "hi.com")
		require.NoError(t, err)
		require.Equal(t, hiStream, stream)
	})

	t.Run("StreamCached", func(t *testing.T) {
		delete(mock.streams, "hi.com")

		stream, err := config.GetStreamConfig(context.TODO(), "hi.com")
		require.NoError(t, err)
		require.Equal(t, hiStream, stream)
	})

	t.Run("StreamMissing", func(t *testing.T) {
		_, err := config.GetStreamConfig(context.TODO(), "hello.com")
		require.Error(t, err)
	})

	t.Run("DevicePassthrough", func(t *testing.T) {
		device, err := config.GetDeviceConfig(context.TODO(), "hiDevice")
		require.NoError(t, err)
		require.Equal(t, hiDevice, device)
	})

	t.Run("DeviceCached", func(t *testing.T) {
		delete(mock.devices, "hi.com")

		device, err := config.GetDeviceConfig(context.TODO(), "hiDevice")
		require.NoError(t, err)
		require.Equal(t, hiDevice, device)
	})

	t.Run("DeviceMissing", func(t *testing.T) {
		_, err := config.GetDeviceConfig(context.TODO(), "helloDevice")
		require.Error(t, err)
	})
}

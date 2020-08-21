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
}

func (m *mockConfigService) GetStreamConfig(ctx context.Context, streamURL string) (data.Stream, error) {
	stream, ok := m.streams[streamURL]
	if !ok {
		return stream, fmt.Errorf("not found")
	}

	return stream, nil
}

func TestCache(t *testing.T) {
	hiStream := data.Stream{
		Secret:      "hi",
		QueryPrefix: "hello",
		Duration:    "1s",
	}
	mock := &mockConfigService{
		streams: map[string]data.Stream{
			"hi.com": hiStream,
		},
	}

	file := os.TempDir() + "/omxplayer-microservice-cache-test.db"
	config, err := New(mock, file)
	require.NoError(t, err)
	defer os.Remove(file)

	t.Run("Passthrough", func(t *testing.T) {
		stream, err := config.GetStreamConfig(context.TODO(), "hi.com")
		require.NoError(t, err)
		require.Equal(t, hiStream, stream)
	})

	t.Run("Cached", func(t *testing.T) {
		delete(mock.streams, "hi.com")

		stream, err := config.GetStreamConfig(context.TODO(), "hi.com")
		require.NoError(t, err)
		require.Equal(t, hiStream, stream)
	})

	t.Run("Missing", func(t *testing.T) {
		_, err := config.GetStreamConfig(context.TODO(), "hello.com")
		require.Error(t, err)
	})
}

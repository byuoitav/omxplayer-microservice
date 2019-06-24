package helpers

import (
	"fmt"
	"time"

	"github.com/godbus/dbus"
)

type OMXPlayer struct {
	Connection *dbus.Conn
	IsReady    bool
}

func (o *OMXPlayer) CanCommand() bool {
	if o.IsReady {
		return true
	}
	//Get playback status
	play, err := GetPlaybackStatus(o.Connection)
	if err == nil && play != "" {
		o.IsReady = true
	}
	return o.IsReady
}

func (o *OMXPlayer) WaitForReady() error {
	for i := 0; !o.IsReady && i < 100; i++ {
		o.CanCommand()
		time.Sleep(50 * time.Millisecond)
	}
	if !o.IsReady {
		return fmt.Errorf("Player media invalid")
	}
	return nil
}

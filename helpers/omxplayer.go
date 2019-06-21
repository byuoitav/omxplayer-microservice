package helpers

import (
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

func (o *OMXPlayer) WaitForReady() {
	for ; !o.IsReady; time.Sleep(100 * time.Millisecond) {
		o.CanCommand()
	}
}

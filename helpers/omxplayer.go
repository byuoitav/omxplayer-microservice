package helpers

import (
	"fmt"
	"time"

	"github.com/godbus/dbus"
)

//OMXPlayer is an interface that reports the status of the dbus connection with the omxplayer
type OMXPlayer struct {
	Connection *dbus.Conn
}

//CanCommand returns a boolean confirming whether it is possible to send commands over dbus
func (o *OMXPlayer) CanCommand() bool {
	_, err := GetPlaybackStatus(o.Connection)
	return err == nil
}

//WaitForReady waits until it is possible to send commands over dbus or until it times out
func (o *OMXPlayer) WaitForReady() error {
	for i := 0; !o.CanCommand() && i < 100; i++ {
		time.Sleep(50 * time.Millisecond)
	}
	if !o.CanCommand() {
		return fmt.Errorf("Player media invalid")
	}
	return nil
}

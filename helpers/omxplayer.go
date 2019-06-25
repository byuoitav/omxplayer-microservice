package helpers

import (
	"fmt"
	"time"

	"github.com/godbus/dbus"
)

type OMXPlayer struct {
	Connection *dbus.Conn
}

func (o *OMXPlayer) CanCommand() bool {
	_, err := GetPlaybackStatus(o.Connection)
	return err == nil
}

func (o *OMXPlayer) WaitForReady() error {
	for i := 0; !o.CanCommand() && i < 100; i++ {
		time.Sleep(50 * time.Millisecond)
	}
	if !o.CanCommand() {
		return fmt.Errorf("Player media invalid")
	}
	return nil
}

package helpers

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/godbus/dbus"
)

type OMXPlayer struct {
	ProcessID  int
	Connection *dbus.Conn
	IsReady    bool
}

func (o *OMXPlayer) IsPlayerRunning() bool {
	process, err := os.FindProcess(o.ProcessID)
	if err != nil {
		return false
	}
	err = process.Signal(syscall.Signal(0))
	return err == nil
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

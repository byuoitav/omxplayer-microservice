package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/godbus/dbus"
)

const (
	envDbusAddress  = "DBUS_SESSION_BUS_ADDRESS"
	envDbusPid      = "DBUS_SESSION_BUS_PID"
	dbusAddressFile = "/tmp/omxplayerdbus.pi"
	dbusIDFile      = "/tmp/omxplayerdbus.pi.pid"
)

//StartOMX ...
func StartOMX(streamURL string) (*OMXPlayer, error) {
	err := runOmxplayer(streamURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to run omxplayer | %s", err.Error())
	}

	err = setEnvironmentVariables()
	if err != nil {
		return nil, fmt.Errorf("Failed to set environment variables | %s", err.Error())
	}

	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to dbus | %s", err.Error())
	}

	omxPlayer := &OMXPlayer{
		Connection: conn,
		IsReady:    false,
	}
	return omxPlayer, err
}

func runOmxplayer(stream string) error {
	cmd := exec.Command("omxplayer", "--avdict", "rtsp_transport:tcp", stream)
	return cmd.Start()
}

func setEnvironmentVariables() error {
	dbusAddress, err := readFile(dbusAddressFile)
	if err != nil {
		return fmt.Errorf("Error when reading dbus address | %s", err.Error())
	}
	err = os.Setenv(envDbusAddress, dbusAddress)
	if err != nil {
		return fmt.Errorf("Error setting dbus address environment variable | %s", err.Error())
	}
	dbusID, err := readFile(dbusIDFile)
	if err != nil {
		return fmt.Errorf("Error when reading dbus id | %s", err.Error())
	}
	err = os.Setenv(envDbusPid, dbusID)
	if err != nil {
		return fmt.Errorf("Error setting dbus id environment variable | %s", err.Error())
	}
	return nil
}

func readFile(path string) (string, error) {
	for i := 0; i < 100; i++ {
		if isExist(path) {
			bytes, err := ioutil.ReadFile(path)
			if err != nil {
				return "", fmt.Errorf("Error when reading file %s | %s", path, err.Error())
			}
			if len(bytes) > 0 {
				return strings.TrimSpace(string(bytes)), err
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	return "", fmt.Errorf("File %s is empty or does not exist", path)
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

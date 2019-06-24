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
	envDbusAddress = "DBUS_SESSION_BUS_ADDRESS"
	envDbusPid     = "DBUS_SESSION_BUS_PID"
)

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
	dbusAddress, err := readFile("/tmp/omxplayerdbus.pi") // Todo: set as constants
	if err != nil {
		return err
	}
	os.Setenv(envDbusAddress, dbusAddress)
	dbusID, err := readFile("/tmp/omxplayerdbus.pi.pid") // Todo: set as constant
	if err != nil {
		return err
	}
	os.Setenv(envDbusPid, dbusID)
	return nil
}

func readFile(path string) (string, error) {
	for i := 0; i < 100; i++ {
		if isExist(path) {
			bytes, err := ioutil.ReadFile(path)
			if err != nil {
				return "", err
			}
			if len(bytes) > 0 {
				return strings.TrimSpace(string(bytes)), err
			}
		}
		time.Sleep(100 * time.Millisecond)
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

package helpers

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/godbus/dbus"
)

const (
	envDbusAddress        = "DBUS_SESSION_BUS_ADDRESS"
	envDbusPid            = "DBUS_SESSION_BUS_PID"
	dbusAddressFilePrefix = "/tmp/omxplayerdbus."
	dbusIDFilePostfix     = ".pid"
	userEnvVar            = "USER"
)

// gross, quick fix! don't call StartOMX from ANYWHERE but handlers.Handlers (the omxMu there must be locked)
var killPrev chan struct{}

//StartOMX starts a new instance of the omxplayer
func StartOMX(streamURL string, args []string) error {
	if err := deleteOmxDbusFiles(); err != nil {
		return fmt.Errorf("unable to delete dbus files: %s", err)
	}

	select {
	case killPrev <- struct{}{}:
	default:
		killPrev = make(chan struct{})
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-killPrev
		cancel()
	}()

	// https://www.raspberrypi.org/documentation/raspbian/applications/omxplayer.md
	// cmd := exec.CommandContext(ctx, "omxplayer", "--display", os.Getenv("OMXPLAYER_DISPLAY"), "--adev", "hdmi", streamURL)
	args = append(args, streamURL)
	cmd := exec.CommandContext(ctx, "omxplayer", args...)
	if err := cmd.Start(); err != nil {
		killPrev <- struct{}{}
		return fmt.Errorf("unable to execute command: %s", err)
	}

	go func() {
		// clean up the process
		if err := cmd.Wait(); err != nil {
			log.L.Warnf("error from omxplayer command: %s", err)
		}
	}()

	return nil
}

// ConnectToDbus establishes a new connection to the omxplayer over dbus
func ConnectToDbus() (*dbus.Conn, error) {
	log.L.Debug("Trying connection to dbus")
	err := setEnvironmentVariables()
	if err != nil {
		log.L.Debugf("Failed to set environment variables | %s", err.Error())
		return nil, fmt.Errorf("Failed to set environment variables | %s", err.Error())
	}

	conn, err := dbus.SessionBusPrivate()
	if err != nil {
		log.L.Debugf("Failed to connect to dbus | %s", err.Error())
		return nil, fmt.Errorf("Failed to connect to dbus | %s", err.Error())
	}

	if err := conn.Auth(nil); err != nil {
		log.L.Debugf("Dbus auth error | %s", err.Error())
		conn.Close()
		conn = nil
		return nil, fmt.Errorf("Failed to connect to dbus, auth error | %s", err.Error())
	}

	if err := conn.Hello(); err != nil {
		log.L.Debugf("Dbus Hello error | %s", err.Error())
		conn.Close()
		conn = nil
	}

	_, err = GetPlaybackStatus(conn)
	return conn, err
}

func deleteOmxDbusFiles() error {
	omxDbusAddressFiles := dbusAddressFilePrefix + "*"
	files, err := filepath.Glob(omxDbusAddressFiles)
	if err != nil {
		return err
	}
	for _, file := range files {
		if err = os.Remove(file); err != nil {
			return err
		}
	}
	return err
}

func setEnvironmentVariables() error {
	userID := os.Getenv(userEnvVar)
	if userID == "" {
		userID = "root"
	}
	log.L.Infof("Environment variable USER: %s", userID)
	dbusAddressFile := dbusAddressFilePrefix + userID
	dbusAddress, err := readFile(dbusAddressFile)
	if err != nil {
		return fmt.Errorf("Error when reading dbus address | %s", err.Error())
	}
	log.L.Debugf("Setenv: DbusAddress: %s", dbusAddress)
	err = os.Setenv(envDbusAddress, dbusAddress)
	if err != nil {
		return fmt.Errorf("Error setting dbus address environment variable | %s", err.Error())
	}
	dbusIDFile := dbusAddressFilePrefix + userID + dbusIDFilePostfix
	dbusID, err := readFile(dbusIDFile)
	if err != nil {
		return fmt.Errorf("Error when reading dbus id | %s", err.Error())
	}
	log.L.Debugf("Setenv: DbusPID: %s", dbusID)
	err = os.Setenv(envDbusPid, dbusID)
	if err != nil {
		return fmt.Errorf("Error setting dbus id environment variable | %s", err.Error())
	}

	log.L.Debugf("Getenv: DbusAddress: %s\nGetenv: DbusPID: %s", os.Getenv(envDbusAddress), os.Getenv(envDbusPid))
	return nil
}

func readFile(path string) (string, error) {
	for i := 0; i < 100; i++ {
		if checkFile(path) {
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

func checkFile(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

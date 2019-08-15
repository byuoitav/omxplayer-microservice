package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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

//StartOMX starts a new instance of the omxplayer and creates an interface through dbus
func StartOMX(streamURL string) (*OMXPlayer, error) {
	log.L.Infof("Removing dbus files")
	err := deleteOmxDbusFiles()
	if err != nil {
		return nil, fmt.Errorf("Failed to remove omxplayer dbus files | %s", err.Error())
	}

	log.L.Infof("Starting omxplayer...")
	err := runOmxplayer(streamURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to run omxplayer | %s", err.Error())
	}

	log.L.Infof("Setting environment variables")
	err = setEnvironmentVariables()
	if err != nil {
		return nil, fmt.Errorf("Failed to set environment variables | %s", err.Error())
	}

	log.L.Infof("Connecting to dbus session")
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to dbus | %s", err.Error())
	}

	omxPlayer := &OMXPlayer{
		Connection: conn,
	}
	return omxPlayer, err
}

func deleteOmxDbusFiles() error {
	omxDbusFiles := dbusAddressFilePrefix + "*"
	return os.Remove(omxDbusFiles)
}

func runOmxplayer(stream string) error {
	cmd := exec.Command("omxplayer", stream)
	return cmd.Start()
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
	err = os.Setenv(envDbusAddress, dbusAddress)
	if err != nil {
		return fmt.Errorf("Error setting dbus address environment variable | %s", err.Error())
	}
	dbusIDFile := dbusAddressFilePrefix + userID + dbusIDFilePostfix
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

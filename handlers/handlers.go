package handlers

import (
	"net/http"
	"net/url"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"
	"github.com/godbus/dbus"

	"github.com/byuoitav/omxplayer-microservice/helpers"

	"github.com/labstack/echo"
)

//PlayStream gets a stream url and attempts to switch the omxplayer output to that stream. If no stream is playing, then a new instance of omxplayer is started.
func PlayStream(ctx echo.Context) error {
	streamURL := ctx.Param("streamURL")
	streamURL, err := url.QueryUnescape(streamURL)
	log.L.Infof("Switching to play stream: %s", streamURL)
	if err != nil {
		log.L.Errorf("Error getting stream URL: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	//Open new connection to dbus
	conn, err := helpers.ConnectToDbus()
	if err != nil {
		log.L.Debug("Can't open dbus connection, starting new stream player")
		err = helpers.StartOMX(streamURL)
		if err != nil {
			log.L.Errorf("Error starting stream player: %s", err.Error())
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		log.L.Infof("Stream player started, playing stream at URL: %s", streamURL)
		return ctx.JSON(http.StatusOK, "Stream player started")
	}

	log.L.Debug("Reconnected to dbus, now switching stream")
	err = switchStream(streamURL, conn)
	if err != nil {
		log.L.Errorf("Error when switching stream: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	log.L.Infof("Stream switched to URL: %s", streamURL)
	return ctx.JSON(http.StatusOK, "Stream switched")
}

func switchStream(streamURL string, conn *dbus.Conn) error {
	if !(checkStream(streamURL, conn)) { // Checks to see if switching to the stream already playing
		err := helpers.SwitchStream(conn, streamURL)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkStream(streamURL string, conn *dbus.Conn) bool {
	currStream, _ := helpers.GetStream(conn)
	return currStream == streamURL // If the streams are the same it returns true, if they are different it returns false
}

//StopStream stops the stream currently running
func StopStream(ctx echo.Context) error {
	log.L.Infof("Stopping stream player...")
	conn, err := helpers.ConnectToDbus()
	if err == nil {
		log.L.Debug("Opened new connection to dbus. Stopping stream player...")
		err := helpers.StopStream(conn)
		if err != nil {
			log.L.Errorf("Error when stopping stream: %s", err.Error())
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		log.L.Infof("Stream player stopped")
		return ctx.JSON(http.StatusOK, "Stream player stopped")
	}
	log.L.Infof("Stream player is not running or is not ready to receive commands")
	return ctx.JSON(http.StatusInternalServerError, "Stream player is not running or is not ready to receive commands")
}

//GetStream returns the url of the stream currently running
func GetStream(ctx echo.Context) error {
	log.L.Infof("Getting current stream URL...")
	conn, err := helpers.ConnectToDbus()
	if err == nil {
		log.L.Debug("Opened new connection to dbus. Getting stream...")
		streamURL, err := helpers.GetStream(conn)
		if err != nil {
			log.L.Errorf("Error when attempting to get current stream: %s", err.Error())
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		log.L.Infof("Getting current stream: %s", streamURL)
		return ctx.JSON(http.StatusOK, status.Input{Input: streamURL})
	}
	log.L.Infof("Stream player is not running or is not ready to receive commands")
	return ctx.JSON(http.StatusInternalServerError, "Stream player is not running or is not ready to receive commands")
}

func checkPlayerStatus(conn *dbus.Conn) error {
	_, err := helpers.GetPlaybackStatus(conn)
	return err
}

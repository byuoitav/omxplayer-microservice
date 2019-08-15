package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"

	"github.com/byuoitav/omxplayer-microservice/helpers"

	"github.com/labstack/echo"
)

var omxPlayer *helpers.OMXPlayer

//PlayStream gets a stream url and attempts to switch the omxplayer output to that stream. If no stream is playing, then a new instance of omxplayer is started.
func PlayStream(ctx echo.Context) error {
	checkPlayerStatus()

	streamURL := ctx.Param("streamURL")
	streamURL, err := url.QueryUnescape(streamURL)
	log.L.Infof("Switching to play stream: %s", streamURL)
	if err != nil {
		log.L.Errorf("Error getting stream URL: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	if omxPlayer == nil {
		err = startNewPlayer(streamURL)
		if err != nil {
			log.L.Errorf("Error starting stream player: %s", err.Error())
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		log.L.Infof("Stream player started, playing stream at URL: %s", streamURL)
		return ctx.JSON(http.StatusOK, "Stream player started")
	}

	err = switchStream(streamURL)
	if err != nil {
		log.L.Errorf("Error when switching stream: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	log.L.Infof("Stream switched to URL: %s", streamURL)
	return ctx.JSON(http.StatusOK, "Stream switched")
}

func startNewPlayer(streamURL string) (err error) {
	omxPlayer, err = helpers.StartOMX(streamURL)
	if err != nil {
		return
	}
	err = omxPlayer.WaitForReady()
	return
}

func switchStream(streamURL string) (err error) {
	if checkStream(streamURL) {
		err = helpers.SwitchStream(omxPlayer.Connection, streamURL)
		if err != nil {
			return
		}
		err = omxPlayer.WaitForReady()
	}
	return
}

func checkStream(streamURL string) bool {
	currStream, _ := helpers.GetStream(omxPlayer.Connection)
	return currStream != streamURL
}

//StopStream stops the stream currently running
func StopStream(ctx echo.Context) error {
	checkPlayerStatus()
	log.L.Infof("Stopping stream player...")
	if omxPlayer != nil {
		err := helpers.StopStream(omxPlayer.Connection)
		if err != nil {
			log.L.Errorf("Error when stopping stream: %s", err.Error())
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		omxPlayer = nil
		log.L.Infof("Stream player stopped")
		return ctx.JSON(http.StatusOK, "Stream player stopped")
	}
	log.L.Infof("Stream player is not running or is not ready to receive commands")
	return ctx.JSON(http.StatusInternalServerError, "Stream player is not running or is not ready to receive commands")
}

//GetStream returns the url of the stream currently running
func GetStream(ctx echo.Context) error {
	checkPlayerStatus()
	log.L.Infof("Getting current stream URL...")
	//Check Player
	if omxPlayer != nil {
		streamURL, err := helpers.GetStream(omxPlayer.Connection)
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

//ChangeVolume ...
func ChangeVolume(ctx echo.Context) error {
	return nil
}

//GetVolume ...
func GetVolume(ctx echo.Context) error {
	checkPlayerStatus()
	if omxPlayer != nil {
		volume, err := helpers.VolumeControl(omxPlayer.Connection)
		if err != nil {
			log.L.Errorf("", err.Error())
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		return ctx.JSON(http.StatusOK, fmt.Sprintf("Current stream volume: %f", volume))
	}
	return ctx.JSON(http.StatusInternalServerError, "Stream player is not running or is not ready to receive commands")
}

//MuteStream ...
func MuteStream(ctx echo.Context) error {
	checkPlayerStatus()
	if omxPlayer != nil {
		err := helpers.Mute(omxPlayer.Connection)
		if err != nil {
			log.L.Errorf("", err.Error())
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		return ctx.JSON(http.StatusOK, "Stream muted")
	}
	return ctx.JSON(http.StatusInternalServerError, "Stream player is not running or is not ready to receive commands")
}

//UnmuteStream ...
func UnmuteStream(ctx echo.Context) error {
	checkPlayerStatus()
	if omxPlayer != nil {
		err := helpers.Unmute(omxPlayer.Connection)
		if err != nil {
			log.L.Errorf("", err.Error())
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		return ctx.JSON(http.StatusOK, "Stream unmuted")
	}
	return ctx.JSON(http.StatusInternalServerError, "Stream player is not running or is not ready to receive commands")
}

func checkPlayerStatus() {
	if omxPlayer != nil && !omxPlayer.CanCommand() {
		omxPlayer = nil
	}
}

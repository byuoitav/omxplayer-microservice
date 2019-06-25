package handlers

import (
	"net/http"
	"net/url"

	"github.com/byuoitav/omxplayer-microservice/helpers"

	"github.com/labstack/echo"
)

var omxPlayer *helpers.OMXPlayer

//PlayStream gets a stream url and attempts to switch the omxplayer output to that stream. If no stream is playing, then a new instance of omxplayer is started.
func PlayStream(ctx echo.Context) error {
	checkPlayerStatus()
	streamURL := ctx.Param("streamURL")
	streamURL, err := url.QueryUnescape(streamURL)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	if omxPlayer == nil {
		err = startNewPlayer(streamURL)
		if err != nil {
			//Todo: Log error
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		return ctx.JSON(http.StatusOK, "Stream player started")
	}
	err = switchStream(streamURL)
	if err != nil {
		//Todo: Log error
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, "Stream switched")
}

func startNewPlayer(streamURL string) (err error) {
	omxPlayer, err = helpers.StartOMX(streamURL)
	if err != nil {
		omxPlayer = nil
		return
	}
	err = omxPlayer.WaitForReady()
	if err != nil {
		omxPlayer = nil
		return
	}
	return
}

func switchStream(streamURL string) (err error) {
	if checkStream(streamURL) {
		err = helpers.SwitchStream(omxPlayer.Connection, streamURL)
		if err != nil {
			return
		}
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
	if omxPlayer != nil && omxPlayer.CanCommand() {
		err := helpers.StopStream(omxPlayer.Connection)
		if err != nil {
			//Todo: Log error
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		omxPlayer = nil
		return ctx.JSON(http.StatusOK, "Stream player stopped")
	}
	return ctx.JSON(http.StatusInternalServerError, "Stream player is not running or is not ready to receive commands")
}

//GetStream returns the url of the stream currently running
func GetStream(ctx echo.Context) error {
	checkPlayerStatus()
	//Check Player
	if omxPlayer != nil && omxPlayer.CanCommand() {
		streamURL, err := helpers.GetStream(omxPlayer.Connection)
		if err != nil {
			//Todo: Log error
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		return ctx.JSON(http.StatusOK, streamURL)
	}
	//Todo: Log error
	return ctx.JSON(http.StatusInternalServerError, "Stream player is not running or is not ready to receive commands")
}

//ChangeVolume ...
func ChangeVolume(ctx echo.Context) error {
	return nil
}

//GetVolume?

//MuteStream ...
func MuteStream(ctx echo.Context) error {
	return nil
}

//UnmuteStream ...
func UnmuteStream(ctx echo.Context) error {
	return nil
}

func checkPlayerStatus() {
	if omxPlayer != nil && !omxPlayer.CanCommand() {
		omxPlayer = nil
	}
}

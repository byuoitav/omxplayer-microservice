package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/byuoitav/omxplayer-microservice/helpers"

	"github.com/labstack/echo"
)

var omxPlayer *helpers.OMXPlayer

//PlayStream ...
func PlayStream(ctx echo.Context) error {
	streamURL := ctx.Param("streamURL")
	streamURL, err := url.QueryUnescape(streamURL)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	if omxPlayer == nil {
		//make a new instance of the player
		omxPlayer, err = helpers.StartOMX(streamURL)
		if err != nil {
			//Log error
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		omxPlayer.WaitForReady()
		return ctx.JSON(http.StatusOK, "Stream player started")
	}
	omxPlayer.WaitForReady()
	err = helpers.SwitchStream(omxPlayer.Connection, streamURL) //Todo: Check if the same stream is already playing
	if err != nil {
		//Log error
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, "Stream switched")
}

//StopStream ...
func StopStream(ctx echo.Context) error {
	if omxPlayer != nil && omxPlayer.CanCommand() {
		err := helpers.StopStream(omxPlayer.Connection)
		if err != nil {
			//Log error
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		omxPlayer = nil
		return ctx.JSON(http.StatusOK, "Stream player stopped")
	}
	return ctx.JSON(http.StatusInternalServerError, fmt.Errorf("Stream player is not running or is not ready to receive commands"))
}

//GetStream ...
func GetStream(ctx echo.Context) error {
	//Check Player
	if omxPlayer != nil && omxPlayer.CanCommand() {
		streamURL, err := helpers.GetStream(omxPlayer.Connection)
		if err != nil {
			//Log error
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		return ctx.JSON(http.StatusOK, streamURL)
	}
	//Log error
	return ctx.JSON(http.StatusInternalServerError, fmt.Errorf("Stream player is not running or is not ready to receive commands"))
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

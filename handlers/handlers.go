package handlers

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/omxplayer-microservice/helpers"

	"github.com/labstack/echo"
)

var omxPlayer *helpers.OMXPlayer

//PlayStream ...
func PlayStream(ctx echo.Context) error {
	streamURL := ctx.Param("streamURL")
	if omxPlayer == nil {
		//make a new instance of the player
		omxPlayer, err := helpers.StartOMX(streamURL)
		if err != nil {
			//Log error
			return ctx.JSON(http.StatusInternalServerError, err)
		}
		omxPlayer.WaitForReady()
		return ctx.JSON(http.StatusOK, nil)
	}
	omxPlayer.WaitForReady()
	err := helpers.SwitchStream(omxPlayer.Connection, streamURL) //Todo: Check if the same stream is already playing
	if err != nil {
		//Log error
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, nil)
}

//StopStream ...
func StopStream(ctx echo.Context) error {
	if omxPlayer != nil && omxPlayer.CanCommand() {
		err := helpers.StopStream(omxPlayer.Connection)
		if err != nil {
			//Log error
			return ctx.JSON(http.StatusInternalServerError, err)
		}
		return ctx.JSON(http.StatusOK, nil)
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
			return ctx.JSON(http.StatusInternalServerError, err)
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

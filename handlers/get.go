package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

func GetTheirStream(ctx echo.Context) error {
	address := ctx.Param("address") // The other pi

	return ctx.JSON(http.StatusOK)
}

func GetMyStream(ctx echo.Context) error {

	return ctx.JSON(http.StatusOK)
}

//player.GetSource()

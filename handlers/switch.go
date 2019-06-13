package handlers

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/common/log"
	"github.com/labstack/echo"
)

func SwitchTheirStream(ctx echo.Context) error {
	address := ctx.Param("address") //The other pi
	streamURL := ctx.Param("streamURL")

	resp, err := http.Get(fmt.Sprintf("http://%s:8032/stream/%s", address, streamURL))
	if err != nil {
		log.L.Errorf("Failed to switch the stream at %s to %s: %s", address, streamURL, err.Error())
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, resp.Body)
}

func SwitchMyStream(ctx echo.Context) error {
	streamURL := ctx.Param("streamURL")

	return ctx.JSON(http.StatusOK)
}

//player.OpenUri

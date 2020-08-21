package handlers

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"crypto/sha256"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"
	_ "github.com/go-kivik/couchdb/v3"
	"github.com/godbus/dbus"

	"github.com/byuoitav/omxplayer-microservice/data"
	"github.com/byuoitav/omxplayer-microservice/helpers"

	"github.com/labstack/echo"
)

type Handlers struct {
	ConfigService     data.ConfigService
	ControlConfigPath string

	omxMu sync.Mutex
}

//PlayStream gets a stream url and attempts to switch the omxplayer output to that stream. If no stream is playing, then a new instance of omxplayer is started.
func (h *Handlers) PlayStream(c echo.Context) error {
	streamURL := c.Param("streamURL")

	if h.ConfigService != nil {
		stream, err := h.ConfigService.GetStreamConfig(c.Request().Context(), streamURL)
		if err == nil && stream.Secret != "" {
			// token is everything after the base url
			token, err := h.generateToken(stream)
			if err != nil {
				log.L.Errorf("error generating secure token: %s", err.Error())
				return c.String(http.StatusInternalServerError, err.Error())
			}

			log.L.Infof("generated secure token: %s\n", token)
			streamURL += token
		}
	}

	streamURL, err := url.QueryUnescape(streamURL)
	if err != nil {
		log.L.Errorf("error getting stream URL: %s", err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// to make sure two go routines don't start omxplayer at the same time
	// and make two instances of it
	h.omxMu.Lock()
	defer h.omxMu.Unlock()

	log.L.Infof("Playing stream %s", streamURL)

	conn, err := helpers.ConnectToDbus()
	if err != nil {
		log.L.Debug("Can't open dbus connection, starting new stream player")

		if err := helpers.StartOMX(streamURL); err != nil {
			log.L.Errorf("Error starting stream player: %s", err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}

		log.L.Infof("Stream player started, playing stream at URL: %s", streamURL)
		return c.String(http.StatusOK, "Stream player started")
	}

	log.L.Debug("Reconnected to dbus, now switching stream")
	if err := h.switchStream(streamURL, conn); err != nil {
		log.L.Errorf("Error when switching stream: %s", err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	log.L.Infof("Successfully started stream %s", streamURL)
	return c.String(http.StatusOK, "Stream switched")
}

func (h *Handlers) switchStream(streamURL string, conn *dbus.Conn) error {
	if !(h.checkStream(streamURL, conn)) { // Checks to see if switching to the stream already playing
		err := helpers.SwitchStream(conn, streamURL)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Handlers) checkStream(streamURL string, conn *dbus.Conn) bool {
	currStream, _ := helpers.GetStream(conn)
	return currStream == streamURL // If the streams are the same it returns true, if they are different it returns false
}

//StopStream stops the stream currently running
func (h *Handlers) StopStream(ctx echo.Context) error {
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
func (h *Handlers) GetStream(ctx echo.Context) error {
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

func (h *Handlers) generateToken(stream data.Stream) (string, error) {
	duration, err := time.ParseDuration(stream.Duration)
	if err != nil {
		return "", fmt.Errorf("unable to parse stream duration: %w", err)
	}

	start := time.Now().UTC()
	end := start.Add(duration)

	url := fmt.Sprintf("%s?%s&%sendtime=%d&%sstarttime=%d", stream.URL, stream.Secret, stream.QueryPrefix, end.Unix(), stream.QueryPrefix, start.Unix())
	input := strings.NewReader(url)
	hash := sha256.New()

	if _, err := io.Copy(hash, input); err != nil {
		return "", fmt.Errorf("unable to copy url to hash")
	}

	finalHash := string(base64.StdEncoding.EncodeToString(hash.Sum(nil)))
	finalHash = strings.ReplaceAll(finalHash, "+", "-")
	finalHash = strings.ReplaceAll(finalHash, "/", "_")

	return fmt.Sprintf("?%sstarttime=%d&%sendtime=%d&%shash=%s", stream.QueryPrefix, start.Unix(), stream.QueryPrefix, end.Unix(), stream.QueryPrefix, finalHash), nil
}

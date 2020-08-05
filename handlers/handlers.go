package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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
	ConfigService data.ConfigService
}

//PlayStream gets a stream url and attempts to switch the omxplayer output to that stream. If no stream is playing, then a new instance of omxplayer is started.
func (h *Handlers) PlayStream(ctx echo.Context) error {
	log.L.Infof("we here\n")
	streamURL := ctx.Param("streamURL")
	// var config data.StreamConfig

	// client, err := kivik.New("couch", fmt.Sprintf("https://%s:%s@%s", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_ADDRESS")))
	// if err != nil {
	// 	log.L.Errorf("error connecting to couch: %s", err.Error())
	// 	return ctx.JSON(http.StatusInternalServerError, err.Error())
	// }

	// db := client.DB(context.TODO(), "stream-configs")
	// if err := db.Get(context.TODO(), "streams").ScanDoc(&config); err != nil {
	// 	log.L.Errorf("error getting stream config doc: %s", err)
	// 	return ctx.JSON(http.StatusInternalServerError, err.Error())
	// }
	// log.L.Infof("stream config: %v", config)

	// if s, ok := config.Streams[streamURL]; ok {
	// 	//generate the hash code and edit the stream url in here
	// 	streamURL, err = h.generateToken(s)
	// 	if err != nil {
	// 		log.L.Errorf("error generating secure token: %s", err.Error())
	// 		return ctx.JSON(http.StatusInternalServerError, err.Error())
	// 	}
	// 	log.L.Infof("generated secure token: %s\n", streamURL)
	// }

	stream, err := h.ConfigService.GetStreamConfig(context.TODO(), streamURL)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	if stream != (data.Stream{}) {
		streamURL, err = h.generateToken(stream)
		if err != nil {
			log.L.Errorf("error generating secure token: %s", err.Error())
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		log.L.Infof("generated secure token: %s\n", streamURL)
	}

	streamURL, err = url.QueryUnescape(streamURL)
	log.L.Infof("Switching to play stream: %s", streamURL)
	if err != nil {
		log.L.Errorf("Error getting stream URL: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	log.L.Infof("final url: %s", streamURL)
	return nil
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
	err = h.switchStream(streamURL, conn)
	if err != nil {
		log.L.Errorf("Error when switching stream: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	log.L.Infof("Stream switched to URL: %s", streamURL)
	return ctx.JSON(http.StatusOK, "Stream switched")
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

func (h *Handlers) generateToken(s data.Stream) (string, error) {
	hash := sha256.New()
	duration, err := time.ParseDuration(s.Duration)
	if err != nil {
		log.L.Errorf("error parsing duration from couch doc: %s", err.Error())
		return "", err
	}
	startTime := time.Now().UTC()
	log.L.Infof("unix time: %v", startTime.Unix())
	endTime := startTime.Add(duration)
	start := startTime.Unix()
	end := endTime.Unix()
	url := fmt.Sprintf("%s?%sendtime=%d&%sstarttime=%d", s.URL, s.QueryPrefix, end, s.QueryPrefix, start)
	input := strings.NewReader(url)
	if _, err := io.Copy(hash, input); err != nil {
		log.L.Errorf("error creating the hash: %s", err.Error())
		return "", err
	}

	finalHash := string(base64.StdEncoding.EncodeToString(hash.Sum(nil)))
	finalHash = strings.ReplaceAll(finalHash, "+", "-")
	finalHash = strings.ReplaceAll(finalHash, "/", "_")
	return fmt.Sprintf("%s?%sstarttime=%d&%sendtime=%d&%shash=%s", s.URL, s.QueryPrefix, start, s.QueryPrefix,
		end, s.QueryPrefix, finalHash), nil
}

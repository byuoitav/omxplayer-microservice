package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/byuoitav/omxplayer-microservice/cache"
	"github.com/byuoitav/omxplayer-microservice/couch"
	"github.com/byuoitav/omxplayer-microservice/handlers"

	"github.com/go-kivik/kivik"

	"github.com/byuoitav/common"
	"github.com/byuoitav/common/log"
)

func main() {
	port := ":8032"
	router := common.NewRouter()

	router.Static("/", "web")

	client, err := kivik.New("couch", fmt.Sprintf("https://%s:%s@%s", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_ADDRESS")))
	if err != nil {
		log.L.Errorf("error connecting to couch: %s", err.Error())
		return
	}

	h := &handlers.Handlers{
		ConfigService: &cache.ConfigService{
			ConfigService: &couch.ConfigService{
				Client:         client,
				StreamConfigDB: "stream-configs",
			},
		},
	}

	router.GET("/control", h.ControlPage)
	router.GET("/stream/:streamURL", h.PlayStream)
	router.GET("/stream/stop", h.StopStream)
	router.GET("/stream", h.GetStream)

	router.PUT("/log-level/:level", log.SetLogLevel)
	router.GET("/log-level", log.GetLogLevel)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}

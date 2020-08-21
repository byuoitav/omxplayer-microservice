package main

import (
	"net/http"
	"net/url"
	"os"

	"github.com/byuoitav/omxplayer-microservice/cache"
	"github.com/byuoitav/omxplayer-microservice/couch"
	"github.com/byuoitav/omxplayer-microservice/handlers"

	_ "github.com/go-kivik/couchdb/v3"
	kivik "github.com/go-kivik/kivik/v3"

	"github.com/byuoitav/common"
	"github.com/byuoitav/common/log"
)

func main() {
	log.SetLevel("info")

	couchURL, err := url.Parse(os.Getenv("DB_ADDRESS"))
	if err != nil {
		log.L.Fatalf("invalid couch address: %s", err)
	}

	couchURL.User = url.UserPassword(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"))
	client, err := kivik.New("couch", couchURL.String())
	if err != nil {
		log.L.Fatalf("error connecting to couch: %s", err)
	}

	couch := &couch.ConfigService{
		Client:         client,
		StreamConfigDB: "stream-configs",
	}

	cache, err := cache.New(couch, os.Getenv("CACHE_DATABASE_LOCATION"))
	if err != nil {
		log.L.Fatalf("unable to build cache: %s", err)
	}

	h := handlers.Handlers{
		ConfigService:     cache,
		ControlConfigPath: os.Getenv("CONTROL_CONFIG_PATH"),
	}

	port := ":8032"
	router := common.NewRouter()

	router.Static("/", "web")
	router.GET("/control", h.ControlPageHandler("./static/control.html"))

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

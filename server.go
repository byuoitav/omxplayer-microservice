package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/byuoitav/omxplayer-microservice/couch"
	"github.com/byuoitav/omxplayer-microservice/handlers"

	_ "github.com/go-kivik/couchdb/v3"
	kivik "github.com/go-kivik/kivik/v3"

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

	config := couch.ConfigService{
		Client:         client,
		StreamConfigDB: "stream-configs",
	}

	h := handlers.Handlers{
		ConfigService: &config,
	}

	// write := router.Group("", auth.AuthorizeRequest("write-state", "room", auth.LookupResourceFromAddress))
	// write.GET("/stream/:streamURL", handlers.PlayStream)
	// write.GET("/stream/stop", handlers.StopStream)
	router.GET("/stream/:streamURL", h.PlayStream)
	router.GET("/stream/stop", h.StopStream)
	// read := router.Group("", auth.AuthorizeRequest("read-state", "room", auth.LookupResourceFromAddress))
	// read.GET("/stream", handlers.GetStream)
	router.GET("/stream", h.GetStream)

	router.PUT("/log-level/:level", log.SetLogLevel)
	router.GET("/log-level", log.GetLogLevel)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}

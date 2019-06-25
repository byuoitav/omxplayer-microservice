package main

import (
	"net/http"

	"github.com/byuoitav/omxplayer-microservice/handlers"

	"github.com/byuoitav/common"
	"github.com/byuoitav/common/log"
)

func main() {
	log.SetLevel("debug")
	port := ":8032"
	router := common.NewRouter()

	// write := router.Group("", auth.AuthorizeRequest("write-state", "room", auth.LookupResourceFromAddress))
	// write.GET("/stream/:streamURL", handlers.PlayStream)
	// write.GET("/stream/stop", handlers.StopStream)

	router.GET("/stream/:streamURL", handlers.PlayStream)
	router.GET("/stream/stop", handlers.StopStream)

	// router.GET("/test", handlers.Test)

	// read := router.Group("", auth.AuthorizeRequest("read-state", "room", auth.LookupResourceFromAddress))
	// read.GET("/stream", handlers.GetStream)

	router.GET("/stream", handlers.GetStream)

	router.PUT("/log-level/:level", log.SetLogLevel)
	router.GET("/log-level", log.GetLogLevel)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}

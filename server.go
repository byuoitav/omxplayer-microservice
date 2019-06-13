package main

import (
	"net/http"

	"github.com/byuoitav/common"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/v2/auth"
	"github.com/byuoitav/omxplayer-microservice/handlers"
)

func main() {
	log.SetLevel("debug")
	port := ":8032"
	router := common.NewRouter()

	write := router.Group("", auth.AuthorizeRequest("write-state", "room", auth.LookupResourceFromAddress))
	write.GET("/:address/stream/:streamURL", handlers.SwitchTheirStream)
	write.GET("/stream/:streamURL", handlers.SwitchMyStream)

	read := router.Group("", auth.AuthorizeRequest("read-state", "room", auth.LookupResourceFromAddress))
	read.GET("/:address/stream", handlers.GetTheirStream)
	read.GET("/stream", handlers.GetMyStream)

	router.PUT("/log-level/:level", log.SetLogLevel)
	router.GET("/log-level", log.GetLogLevel)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}

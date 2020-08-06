package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/byuoitav/omxplayer-microservice/cache"
	"github.com/byuoitav/omxplayer-microservice/couch"
	"github.com/byuoitav/omxplayer-microservice/handlers"

	_ "github.com/go-kivik/couchdb/v3"
	kivik "github.com/go-kivik/kivik/v3"
	bolt "go.etcd.io/bbolt"

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

	c := couch.ConfigService{
		Client:         client,
		StreamConfigDB: "stream-configs",
	}

	dbLoc := os.Getenv("CACHE_DATABASE_LOCATION")
	db, err := bolt.Open(dbLoc, 0600, nil)
	if err != nil {
		log.L.Errorf("error creating cache: %s", err.Error())
		return
	}

	err = db.Update(func(tx *bolt.Tx) error {
		log.L.Debugf("Checking if Stream Bucket Exists")
		_, err := tx.CreateBucketIfNotExists([]byte("STREAMS"))
		if err != nil {
			return fmt.Errorf("error creating stream bucket: %s", err.Error())
		}

		return nil
	})
	if err != nil {
		log.L.Errorf("could not create db bucket: %s", err.Error())
		return
	}

	config := cache.ConfigService{
		ConfigService: &c,
		DB:            db,
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

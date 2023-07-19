package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"project/internal/services/nats/subscriber"
	"project/internal/storage/psql"
	rediscache "project/internal/storage/redis"
	"project/internal/web/handlers"
	"time"

	"github.com/nats-io/nats.go"
)

type Subscriber interface {
	Subscribe(subject string, handler func(msg *nats.Msg)) error
	ServSubscription(subj string) error
}

type cfg struct {
	cacheAddr       string
	cachePassword   string
	cacheDB         int
	cacheExparation time.Duration
	databaseConnStr string
}

func main() {
	cfg := cfg{
		cacheAddr:       "localhost:6379",
		cacheDB:         0,
		cacheExparation: 10 * time.Hour,
		databaseConnStr: "postgres://nikolay:pass@localhost:5432/wb_l0?sslmode=disable",
	}

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Panic(err)
	}

	db := psql.NewPostgresConn(cfg.databaseConnStr)
	cache := rediscache.NewRedisCache(cfg.cacheDB, cfg.cacheAddr, cfg.cachePassword, cfg.cacheExparation)
	subscriber, err := subscriber.NewSubscriber(nc, db, cache)
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()
	var subj = "test"
	args := flag.Args()
	if len(args) != 0 {
		subj = args[0]
	}

	err = subscriber.Subscribe(subj)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Web Server running at http://localhost:8080")
	http.HandleFunc("/", handlers.HandleMain(cache))
	http.ListenAndServe(":8080", nil)

	// err = subscriber.ServSubscription()
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }
}

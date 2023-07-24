package main

import (
	"flag"
	"log"
	"project/internal/services/nats/subscriber"
	"project/internal/storage/inmem"
	"project/internal/storage/psql"
	"project/internal/web/server"
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
	cache := inmem.NewCache(cfg.cacheExparation, cfg.cacheExparation*2)
	err = cache.Restore(db)
	cache.Lookup()
	if err != nil {
		log.Fatal(err)
	}

	subscriber, err := subscriber.NewSubscriber(nc, db, cache)
	if err != nil {
		log.Fatal(err)
	}

	server := server.NewServer(cache)

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
	defer subscriber.Unsubscribe()

	server.RunServer(cache)
}

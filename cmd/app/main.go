package main

import (
	"log"
	"os"
	"project/internal/models"
	"project/internal/services/nats-stan/subscriber"
	"project/internal/storage/psql"
	rediscache "project/internal/storage/redis"
	"time"

	"github.com/nats-io/nats.go"
)

type Subscriber interface {
	Subscribe(subject string, handler func(msg *nats.Msg)) error
	ServSubscription(subj string) error
}

type db interface {
	SaveDataFromOrder(order models.Order) error
	GetOrdersInLast(time.Duration) ([]models.Order, error)
}

type cache interface {
	GetOrder(orederID string) (models.Order, error)
	SetOrder(order models.Order) error
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

	nc, _ := nats.Connect(nats.DefaultURL)

	db := psql.NewPostgresConn(cfg.databaseConnStr)
	cache := rediscache.NewRedisCache(cfg.cacheDB, cfg.cacheAddr, cfg.cachePassword, cfg.cacheExparation)
	subscriber, err := subscriber.NewSubscriber(nc, db, cache)
	if err != nil {
		log.Fatal(err)
	}

	err = subscriber.Subscribe("order")
	if err != nil {
		log.Fatal(err)
	}

	err = subscriber.ServSubscription()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

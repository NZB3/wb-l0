package main

import (
	"log"
	"os"
	"project/internal/models"
	"project/internal/service/nats/subscriber"
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

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	var db db
	var cache cache
	subscriber, err := subscriber.NewSubscriber(nc, db, cache)
	if err != nil {
		log.Fatal(err)
	}

	err = subscriber.Subscribe("foo")
	if err != nil {
		log.Fatal(err)
	}

	err = subscriber.ServSubscription()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

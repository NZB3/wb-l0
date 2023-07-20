package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"project/internal/models"
	"project/internal/services/nats-stan/publisher"

	"github.com/nats-io/nats.go"
)

type Publisher interface {
	Publish(data []byte, subject string) error
}

func main() {
	var err error

	flag.Parse()
	args := flag.Args()

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Panic(err)
	}
	defer nc.Close()
	publisher, err := publisher.NewPublisher(nc)
	if err != nil {
		log.Fatal(err)
	}

	subject := "test"
	data := []byte("Hello world!")

	if len(args) != 0 {
		subject = args[0]
		data, err = ioutil.ReadFile("data.json")
		if err != nil {
			log.Fatal(err)
		}

		orders := make([]models.Order, 1)
		json.Unmarshal(data, &orders)
		for _, order := range orders {
			data, err := json.Marshal(order)
			if err != nil {
				log.Fatal(err)
			}
			if err := publisher.Publish(data, subject); err != nil {
				log.Fatal(err)
			}
		}
	}

	if err := publisher.Publish(data, subject); err != nil {
		log.Fatal(err)
	}
}

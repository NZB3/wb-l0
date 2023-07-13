package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"project/internal/storage/models"
	"project/pkg/cache"
	"project/pkg/database"
	"time"

	"github.com/lib/pq"
)

type cfg struct {
	cacheAddr       string
	cachePassword   string
	cacheDB         int
	cacheExparation time.Duration
	databaseConnStr string
	natsID          string
	natsURL         string
}

func main() {
	cfg := cfg{
		cacheAddr:       "localhost:6379",
		cachePassword:   "",
		cacheDB:         0,
		cacheExparation: 10 * time.Hour,
		databaseConnStr: "postgres://nikolay:pass@localhost:5432/wb_l0?sslmode=disable",
		natsID:          "NBSTRQF62HII25U74VYAIEMTIHOUW5Z6S3BPX4WTVDTY2NGPPU2BUXM6",
		natsURL:         "nats://localhost:8222",
	}

	db, err := database.NewConnection(cfg.databaseConnStr)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	cache := cache.New(cfg.cacheAddr, cfg.cachePassword, cfg.cacheDB, cfg.cacheExparation)

	data, err := ioutil.ReadFile("data.json")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	var orders []models.Order
	json.Unmarshal([]byte(data), &orders)
	fmt.Println(orders)

	for _, order := range orders {
		err = db.SaveDataFromOrder(order)
		if err != nil {
			log.Println(err)
		}
	}

	fmt.Println()

	for _, order := range orders {
		err = cache.SetOrder(order)
		if err != nil {
			log.Println(err)
		}
	}

	fmt.Println()

	order, err := cache.GetOrder("b563feb7b2b84b6test")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(order)

	fmt.Println()

	err = cache.Clear()
	if err != nil {
		log.Println(err)
	}
	log.Println("cache cleared")
	fmt.Println()

	log.Println("start restore from db")
	err = cache.RestoreFrom(db)
	if err != nil {
		log.Println(err)
	}

	order, err = cache.GetOrder("b563feb7b2b84b6test1")
	if err != nil {
		log.Println(err)
	}

	fmt.Println(order)
	_ = pq.ErrSSLNotSupported
	_ = cache
	_ = db
}

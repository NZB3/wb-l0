package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"project/internal/storage/cache"
	"project/internal/storage/database"
	"project/internal/storage/models"
	"time"

	"github.com/lib/pq"
)

type cfg struct {
	cacheAddr       string
	cachePassword   string
	cacheDB         int
	cacheExparation time.Duration
	databaseConnStr string
}

const data = `
{
  "order_uid": "b563feb7b2b84b6test",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b6test",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    },
	{
		"chrt_id": 9934921,
		"track_number": "WBILMTESTTRACK",
		"price": 453,
		"rid": "ab4219087a764ae1btest",
		"name": "shoes",
		"sale": 33,
		"size": "1",
		"total_price": 200,
		"nm_id": 2389213,
		"brand": "Nike",
		"status": 202
	  }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}
`

func main() {
	cfg := cfg{
		cacheAddr:       "localhost:6379",
		cachePassword:   "",
		cacheDB:         0,
		cacheExparation: 10 * time.Hour,
		databaseConnStr: "postgres://nikolay:pass@localhost:5432/wb_l0?sslmode=disable",
	}

	db, err := database.NewConnection(cfg.databaseConnStr)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	cache := cache.New(cfg.cacheAddr, cfg.cachePassword, cfg.cacheDB, cfg.cacheExparation)

	order := models.Order{}
	json.Unmarshal([]byte(data), &order)

	err = db.SaveDataFromOrder(order)
	if err != nil {
		log.Println(err)
	}

	fmt.Println()

	err = cache.SetOrder(order)
	if err != nil {
		log.Println(err)
	}

	fmt.Println()

	order, err = cache.GetOrder(order.OrderUID)
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

	order, err = cache.GetOrder(order.OrderUID)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(order)
	_ = pq.ErrSSLNotSupported
	_ = cache
	_ = db
}

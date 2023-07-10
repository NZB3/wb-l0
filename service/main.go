package main

import (
	"project/internal/storage/cache"
	"project/internal/storage/database"

	"github.com/lib/pq"
)

type cfg struct {
	cacheAddr       string
	cachePassword   string
	cacheDB         int
	databaseConnStr string
}

func main() {
	cfg := cfg{
		cacheAddr:       "localhost:6379",
		cachePassword:   "",
		cacheDB:         0,
		databaseConnStr: "postgres://nikolay:pass@localhost:5432/wb_l0?sslmode=disable",
	}

	cache := cache.New(cfg.cacheAddr, cfg.cachePassword, cfg.cacheDB)
	db := database.NewConnection(cfg.databaseConnStr)

	_ = pq.ErrSSLNotSupported
	_ = cache
	_ = db
}

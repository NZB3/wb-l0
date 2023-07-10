package cache

import (
	"github.com/go-redis/redis"
)

type cache struct {
	client *redis.Client
}

var instance *cache = nil

func New(addr, password string, db int) *cache {
	if instance == nil {
		instance = &cache{
			client: new(addr, password, db),
		}
	}

	return instance
}

func new(addr, password string, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return client
}

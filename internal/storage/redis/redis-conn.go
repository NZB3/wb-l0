package rediscache

import (
	"log"
	"project/internal/storage"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type redisCache struct {
	client         *redis.Client
	exparationTime time.Duration
}

var instance *redisCache = nil
var once sync.Once

func NewRedisCache(db int, addr, password string) *redisCache {
	once.Do(func() {
		instance = &redisCache{
			client: newRedisClient(db, addr, password),
		}
	})
	return instance
}

func (r *redisCache) CheckCache() error {
	if r.client == nil {
		return storage.ErrCacheNotExists
	}
	return nil
}

func newRedisClient(db int, addr, password string) *redis.Client {
	op := "storage.redis.newRedisClient"
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Panicf("%s: %s", op, err)
	}
	return client
}

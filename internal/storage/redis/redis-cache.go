package rediscache

import (
	"project/internal/storage"

	"github.com/go-redis/redis"
)

type redisCache struct {
	client *redis.Client
}

func (r *redisCache) CheckCache() error {
	if r.client == nil {
		return storage.ErrCacheNotExists
	}
	return nil
}

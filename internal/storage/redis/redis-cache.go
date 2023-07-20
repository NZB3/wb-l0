package rediscache

import (
	"encoding/json"
	"fmt"
	"project/internal/models"
	"time"
)

type db interface {
	GetOrdersInLast(time.Duration) ([]models.Order, error)
}

func (rc *redisCache) GetExparationTime() time.Duration {
	return rc.exparationTime
}

func (rc *redisCache) Pong() {
	pong, err := rc.client.Ping().Result()
	fmt.Println(pong, err)
}

func (rc *redisCache) GetOrder(orederUID string) ([]byte, error) {
	const op = "storage.redis-cache.GetOrder"

	jsonOrder, err := rc.client.Get(orederUID).Bytes()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	return jsonOrder, nil
}

func (rc *redisCache) SetOrder(order models.Order) error {
	const op = "storage.redis-cache.SetOrder"

	orderJSON, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	err = rc.client.Set(order.OrderUID, orderJSON, rc.exparationTime).Err()
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}
	return nil
}

func (rc *redisCache) RestoreFrom(db db) error {
	const op = "storage.redis-cache.RestoreFrom"

	orders, err := db.GetOrdersInLast(rc.exparationTime)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}
	for _, order := range orders {
		err = rc.SetOrder(order)
		if err != nil {
			return fmt.Errorf("%s: %s", op, err)
		}
	}
	return nil
}

func (rc *redisCache) Clear() error {
	const op = "storage.redis-cache.Clear"

	err := rc.client.FlushDB().Err()
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	return nil
}

package cache

import (
	"encoding/json"
	"fmt"
	"project/internal/storage/database"
	"project/internal/storage/models"
	"time"
)

type Cache interface {
	GetOrder(orederID string) (models.Order, error)
	SetOrder(order models.Order) error
	GetExparationTime() time.Duration
	Clear() error
	Pong() (string, error)
}

func (cache *cache) GetExparationTime() time.Duration {
	return cache.exparationTime
}

func (cache *cache) Pong() {
	pong, err := cache.client.Ping().Result()
	fmt.Println(pong, err)
}

func (cache *cache) GetOrder(orederUID string) (models.Order, error) {
	const op = "cache.GetOrder"

	cachedOrder, err := cache.client.Get(orederUID).Bytes()
	if err != nil {
		return models.Order{}, fmt.Errorf("%s: %s", op, err)
	}

	order := models.Order{}
	err = json.Unmarshal(cachedOrder, &order)
	if err != nil {
		return models.Order{}, fmt.Errorf("%s: %s", op, err)
	}

	return order, nil
}

func (cache *cache) SetOrder(order models.Order) error {
	const op = "cache.SetOrder"

	orderJSON, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	err = cache.client.Set(order.OrderUID, orderJSON, cache.exparationTime).Err()
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}
	return nil
}

func (cache *cache) RestoreFrom(db database.DB) error {
	const op = "cache.Restore"

	orders, err := db.GetOrdersInLast(cache.exparationTime)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}
	for _, order := range orders {
		err = cache.SetOrder(order)
		if err != nil {
			return fmt.Errorf("%s: %s", op, err)
		}
	}
	return nil
}

func (cache *cache) Clear() error {
	const op = "cache.Clear"

	err := cache.client.FlushDB().Err()
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	return nil
}

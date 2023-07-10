package cache

import (
	"encoding/json"
	"fmt"
	"project/internal/storage/models"
	"time"
)

type Cach interface {
	GetOrder(orederID string) (models.Order, error)
	SetOrder(order models.Order) error
	Pong() (string, error)
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

	err = cache.client.Set(order.OrderUID, orderJSON, 5*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	return nil
}

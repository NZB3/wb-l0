package inmem

import (
	"encoding/json"
	"fmt"
	"log"
	"project/internal/models"
	"project/internal/storage"
	"time"

	"github.com/patrickmn/go-cache"
)

type db interface {
	GetOrdersInLast(duration time.Duration) ([]models.Order, error)
}

type inmemCache struct {
	orders         *cache.Cache
	expirationTime time.Duration
}

func NewCache(expirationTime, cleanUpTime time.Duration) *inmemCache {
	cache := cache.New(expirationTime, cleanUpTime)
	return &inmemCache{
		orders:         cache,
		expirationTime: expirationTime,
	}
}

func (c *inmemCache) GetOrder(orederUID string) ([]byte, error) {
	const op = "storage.inmem-cache.GetOrder"

	order, ok := c.orders.Get(orederUID)
	if !ok {
		return nil, fmt.Errorf("%s: %s", op, storage.ErrOrderNotFound)
	}

	return order.([]byte), nil
}

func (c *inmemCache) SetOrder(order models.Order) error {
	const op = "storage.inmem-cache.SetOrder"
	jsonOrder, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	err = c.orders.Add(order.OrderUID, jsonOrder, c.expirationTime)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	return nil
}

func (c *inmemCache) Restore(db db) error {
	const op = "storage.inmem-cache.RestoreFrom"

	orders, err := db.GetOrdersInLast(c.expirationTime)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	if len(orders) == 0 {
		log.Printf("%s: no orders to restore", op)
		return nil
	}

	for _, order := range orders {
		err = c.SetOrder(order)
		if err != nil {
			return fmt.Errorf("%s: %s", op, err)
		}
	}

	return nil
}

func (c *inmemCache) Clear() error {
	c.orders.Flush()

	return nil
}

func (c *inmemCache) Lookup() {
	fmt.Println(c.orders.Items())
}

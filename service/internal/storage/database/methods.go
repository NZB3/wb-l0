package database

import (
	"project/internal/storage/models"
)

type DB interface {
	SaveOrder(order models.Order) error
	GetOrder(orderUID string) (models.Order, error)
	GetAllOrders() ([]models.Order, error)
}

func (db *db) SaveOrder(order models.Order) error {
	const op = "db.SaveOrder"

}

func (db *db) GetAllOrders() ([]models.Order, error) {
	op := "db.GetAllOrders"
}

func (db *db) GetOrder(orderUID string) (models.Order, error) {
	op := "db.GetOrder"
}

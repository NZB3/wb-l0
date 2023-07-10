package database

import (
	"fmt"
	"project/internal/storage/models"
)

type DB interface {
	SaveOrder(order models.Order) error
}

func (db *db) SaveOrder(order models.Order) error {
	const op = "db.SaveOrder"

	tx, err := db.connection.Begin()
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}
	defer tx.Rollback()

	// Insert data into the orders table
	_, err = tx.Exec("INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey, order.SMID, order.DateCreated, order.OofShard)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %s", op, err)
	}

	// Insert data into the delivery table
	_, err = tx.Exec("INSERT INTO delivery (order_uid, recipient_name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %s", op, err)
	}

	// Insert data into the payment table
	_, err = tx.Exec("INSERT INTO payments (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %s", op, err)
	}

	// Insert data into the items table
	for _, item := range order.Items {
		_, err = tx.Exec("INSERT INTO items (order_uid, chrt_id, track_number, price, rid, chrt_name, sale, chrt_size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("%s: %s", op, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	return nil
}

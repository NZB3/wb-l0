package psql

import (
	"database/sql"
	"fmt"
	"project/internal/models"
	"project/internal/storage"
	"time"
)

func (pq *psql) SaveDataFromOrder(order models.Order) error {
	const op = "storage.psql.SaveDataFromOrder"

	tx, err := pq.conn.Begin()
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}
	defer tx.Rollback()

	err = saveOrder(order, tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %s", op, err)
	}

	err = saveDelivery(order, tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %s", op, err)
	}

	err = savePayment(order, tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %s", op, err)
	}

	err = saveItems(order, tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %s", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	return nil
}

func (pq *psql) GetOrdersInLast(duration time.Duration) ([]models.Order, error) {
	const op = "db.GetOrdersInLast"

	startTime := time.Now().Add(duration * -1)
	rows, err := pq.conn.Query(`
		SELECT o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
		d.recipient_name, d.phone, d.zip, d.city, d.address, d.region, d.email,
		p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
		FROM orders o
		LEFT JOIN delivery d ON o.order_uid = d.order_uid
		LEFT JOIN payments p ON o.order_uid = p.order_uid
		WHERE o.date_created >= $1;
	`, startTime)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var orders []models.Order

	for rows.Next() {
		var order models.Order
		delivery := &order.Delivery
		payment := &order.Payment
		items := &order.Items

		err = rows.Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SMID, &order.DateCreated.Time, &order.OofShard,
			&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email,
			&payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider, &payment.Amount, &payment.PaymentDT, &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		rows, err := pq.conn.Query(`
			SELECT i.chrt_id, i.track_number, i.price, i.rid, i.chrt_name, i.sale, i.chrt_size, i.total_price, i.nm_id, i.brand, i.status
			FROM items i WHERE i.order_uid = $1;
		`, order.OrderUID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		defer rows.Close()

		for rows.Next() {
			var item models.Item
			err = rows.Scan(
				&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name,
				&item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
			)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}

			*items = append(*items, item)
		}

		orders = append(orders, order)

	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return orders, nil
}

func (p *psql) CheckDB() error {
	if p.conn == nil {
		return storage.ErrDBNotExists
	}
	return nil
}

func saveDelivery(order models.Order, tx *sql.Tx) error {
	const op = "db.saveDelivery"

	_, err := tx.Exec("INSERT INTO delivery (order_uid, recipient_name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	return nil
}

func savePayment(order models.Order, tx *sql.Tx) error {
	const op = "db.savePayment"

	_, err := tx.Exec("INSERT INTO payments (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	return nil
}

func saveItems(order models.Order, tx *sql.Tx) error {
	const op = "db.saveItems"

	for _, item := range order.Items {
		_, err := tx.Exec("INSERT INTO items (order_uid, chrt_id, track_number, price, rid, chrt_name, sale, chrt_size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return fmt.Errorf("%s: %s", op, err)
		}
	}

	return nil
}

func saveOrder(order models.Order, tx *sql.Tx) error {
	const op = "db.saveOrder"
	_, err := tx.Exec("INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey, order.SMID, order.DateCreated.Time, order.OofShard)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	return nil
}

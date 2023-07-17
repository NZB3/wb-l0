package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Order struct {
	OrderUID          string      `json:"order_uid"`
	TrackNumber       string      `json:"track_number"`
	Entry             string      `json:"entry"`
	Delivery          Delivery    `json:"delivery"`
	Payment           Payment     `json:"payment"`
	Items             []Item      `json:"items"`
	Locale            string      `json:"locale"`
	InternalSignature string      `json:"internal_signature"`
	CustomerID        string      `json:"customer_id"`
	DeliveryService   string      `json:"delivery_service"`
	ShardKey          string      `json:"shardkey"`
	SMID              int         `json:"sm_id"`
	DateCreated       dateCreated `json:"date_created"`
	OofShard          string      `json:"oof_shard"`
}
type dateCreated struct {
	time.Time
}

func (d *dateCreated) UnmarshalJSON(data []byte) error {
	op := "models.dateCreated.UnmarshalJSON"

	var t time.Time
	err := json.Unmarshal(data, &t)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if t.IsZero() {
		t = time.Now()
	}
	d.Time = t
	return nil
}

func (d *dateCreated) MarshalJSON() ([]byte, error) {
	op := "models.dateCreated.MarshalJSON"

	data, err := json.Marshal(d.Time)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return data, nil
}

func (o *Order) Unmarshal(data []byte) error {
	op := "models.Order.Unmarshal"

	err := json.Unmarshal(data, o)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

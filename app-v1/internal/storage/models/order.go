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
	DateCreated       TimeCreated `json:"date_created"`
	OofShard          string      `json:"oof_shard"`
}
type TimeCreated struct {
	Time time.Time
}

func (t *TimeCreated) UnmarshalJSON(data []byte) error {
	op := "TimeCreated.UnmarshalJSON"

	if string(data) == "null" {
		*t = TimeCreated{time.Now()}
	} else {
		err := json.Unmarshal(data, &t.Time)
		if err != nil {
			return fmt.Errorf("%s: %s", op, err)
		}
	}
	return nil
}

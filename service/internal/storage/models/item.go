package models

import "github.com/Rhymond/go-money"

type Item struct {
	ChrtID      int         `json:"chrt_id"`
	TrackNumber string      `json:"track_number"`
	Price       money.Money `json:"price"`
	RID         string      `json:"rid"`
	Name        string      `json:"name"`
	Sale        int         `json:"sale"`
	Size        string      `json:"size"`
	TotalPrice  money.Money `json:"total_price"`
	NmID        int         `json:"nm_id"`
	Brand       string      `json:"brand"`
	Status      int         `json:"status"`
}

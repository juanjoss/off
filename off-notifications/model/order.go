package model

import "time"

type Order struct {
	Id        int       `json:"id" db:"id"`
	SsdId     int       `json:"ssd_id" db:"ssd_id"`
	ProductId string    `json:"product_id" db:"product_id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	Quantity  int       `json:"quantity" db:"quantity"`
	Status    string    `json:"status" db:"status"`
}

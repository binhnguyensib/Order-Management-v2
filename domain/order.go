package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Order struct {
	Id          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	CustomerId  string        `json:"customer_id" `
	ProductIds  []string      `json:"product_ids"`
	TotalAmount float64       `json:"total_amount"`
	CreatedAt   time.Time     `json:"created_at"`
}

type OrderRequest struct {
	CustomerId  string   `json:"customer_id"`
	ProductIds  []string `json:"product_ids"`
	TotalAmount float64  `json:"total_amount"`
}

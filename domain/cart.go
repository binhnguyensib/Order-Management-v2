package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type Cart struct {
	Id         bson.ObjectID `json:"id" bson:"_id,omitempty"`
	CustomerID string        `json:"customer_id" bson:"customer_id"`
	Items      []*CartItem   `json:"items" bson:"items"`
	TotalItems int           `json:"total_items" bson:"total_items"`
	TotalPrice float64       `json:"total_price" bson:"total_price"`
}

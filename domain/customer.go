package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type Customer struct {
	Id    bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Name  string        `json:"name"`
	Email string        `json:"email"`
	Phone string        `json:"phone"`
}

type CustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

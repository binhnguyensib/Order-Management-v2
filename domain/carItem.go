package domain

type CartItem struct {
	ProductID    string  `json:"product_id" bson:"product_id"`
	ProductName  string  `json:"product_name" bson:"product_name"`
	ProductPrice float64 `json:"product_price" bson:"product_price"`
	Quantity     int     `json:"quantity" bson:"quantity"`
	Subtotal     float64 `json:"subtotal" bson:"subtotal"`
}

type CartItemRequest struct {
	ProductID   string `json:"product_id" bson:"product_id"`
	ProductName string `json:"product_name" bson:"product_name"`
	Quantity    int    `json:"quantity" bson:"quantity"`
}

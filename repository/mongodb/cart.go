package mongodb

import (
	"context"
	"fmt"
	"intern-project-v2/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type cartRepositoryImpl struct {
	conn *mongo.Database
}

func NewCartRepository(db *mongo.Database) domain.CartRepository {
	return &cartRepositoryImpl{
		conn: db,
	}
}

func (cr *cartRepositoryImpl) AddToCart(ctx context.Context, customerID string, item *domain.CartItem) (*domain.Cart, error) {

	collection := cr.conn.Collection("carts")
	var existingCart domain.Cart
	err := collection.FindOne(ctx, bson.M{"customer_id": customerID}).Decode(&existingCart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			newCart := &domain.Cart{
				CustomerID: customerID,
				Items:      []*domain.CartItem{item},
				TotalItems: item.Quantity,
				TotalPrice: item.ProductPrice * float64(item.Quantity),
			}
			result, err := collection.InsertOne(ctx, newCart)
			if err != nil {
				return nil, fmt.Errorf("failed to create new cart: %v", err)
			}

			if insertedId, ok := result.InsertedID.(bson.ObjectID); ok {
				newCart.Id = insertedId
			}
			return newCart, nil
		}
		return nil, fmt.Errorf("failed to find cart for customer %s: %v", customerID, err)
	}
	found := false
	for i, existingItem := range existingCart.Items {
		if existingItem.ProductID == item.ProductID {
			existingCart.Items[i].Quantity += item.Quantity
			existingCart.TotalItems += item.Quantity
			existingCart.TotalPrice += item.ProductPrice * float64(item.Quantity)
			found = true
			break
		}
	}
	if !found {
		existingCart.Items = append(existingCart.Items, item)
	}
	cr.recalCartTotals(&existingCart)
	_, err = collection.UpdateOne(ctx, bson.M{"_id": existingCart.Id}, bson.M{"$set": existingCart})
	if err != nil {
		return nil, fmt.Errorf("failed to update cart: %v", err)
	}
	return &existingCart, nil

}

func (cr *cartRepositoryImpl) GetCartByCustomerId(ctx context.Context, customerID string) (*domain.Cart, error) {
	collection := cr.conn.Collection("carts")
	var cart domain.Cart
	err := collection.FindOne(ctx, bson.M{"customer_id": customerID}).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("customer's cart is empty")
		}
		return nil, fmt.Errorf("failed to find cart for customer %s: %v", customerID, err)
	}
	return &cart, nil
}

func (cr *cartRepositoryImpl) UpdateCartItem(ctx context.Context, customerID string, item *domain.CartItem) (*domain.Cart, error) {
	collection := cr.conn.Collection("carts")
	var existingCart domain.Cart
	err := collection.FindOne(ctx, bson.M{"customer_id": customerID}).Decode(&existingCart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("customer's cart is empty")
		}
		return nil, fmt.Errorf("failed to find cart for customer %s: %v", customerID, err)
	}

	found := false
	for i, existingItem := range existingCart.Items {
		if existingItem.ProductID == item.ProductID {
			existingCart.Items[i].Quantity = item.Quantity
			existingCart.Items[i].Subtotal = item.ProductPrice * float64(item.Quantity)
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("item with product ID %s not found in cart", item.ProductID)
	}

	cr.recalCartTotals(&existingCart)
	_, err = collection.UpdateOne(ctx, bson.M{"_id": existingCart.Id}, bson.M{"$set": existingCart})
	if err != nil {
		return nil, fmt.Errorf("failed to update cart: %v", err)
	}
	return &existingCart, nil
}

func (cr *cartRepositoryImpl) RemoveCartItem(ctx context.Context, customerID string, productID string) (*domain.Cart, error) {
	collection := cr.conn.Collection("carts")
	var existingCart domain.Cart
	err := collection.FindOne(ctx, bson.M{"customer_id": customerID}).Decode(&existingCart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("customer's cart is empty")
		}
		return nil, fmt.Errorf("failed to find cart for customer %s: %v", customerID, err)
	}

	var updatedItems []*domain.CartItem
	for _, item := range existingCart.Items {
		if item.ProductID != productID {
			updatedItems = append(updatedItems, item)
		}
	}

	if len(updatedItems) == len(existingCart.Items) {
		return nil, fmt.Errorf("item with product ID %s not found in cart", productID)
	}

	existingCart.Items = updatedItems
	cr.recalCartTotals(&existingCart)

	_, err = collection.UpdateOne(ctx, bson.M{"_id": existingCart.Id}, bson.M{"$set": existingCart})
	if err != nil {
		return nil, fmt.Errorf("failed to update cart: %v", err)
	}
	return &existingCart, nil
}

func (cr *cartRepositoryImpl) ClearCart(ctx context.Context, customerID string) error {
	collection := cr.conn.Collection("carts")
	result, err := collection.DeleteOne(ctx, bson.M{"customer_id": customerID})
	if err != nil {
		return fmt.Errorf("failed to clear cart for customer %s: %v", customerID, err)
	}
	if result.DeletedCount == 0 {
		fmt.Printf("customer's cart has already empty")
	}
	return nil
}

func (cr *cartRepositoryImpl) recalCartTotals(cart *domain.Cart) {
	totalItems := 0
	totalPrice := 0.0
	for _, item := range cart.Items {
		totalItems += item.Quantity
		totalPrice += item.ProductPrice * float64(item.Quantity)
	}
	cart.TotalItems = totalItems
	cart.TotalPrice = totalPrice
}

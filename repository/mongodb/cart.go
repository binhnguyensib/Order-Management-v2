package mongodb

import (
	"context"
	"intern-project-v2/domain"
	"intern-project-v2/logger"

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
	logger.Info("Adding item to cart", "customerID", customerID, "item", item)
	collection := cr.conn.Collection("carts")
	var existingCart domain.Cart
	err := collection.FindOne(ctx, bson.M{"customer_id": customerID}).Decode(&existingCart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Info("Creating new cart for customer",
				"customer_id", customerID,
				"reason", "cart_not_found",
			)
			newCart := &domain.Cart{
				CustomerID: customerID,
				Items:      []*domain.CartItem{item},
				TotalItems: item.Quantity,
				TotalPrice: item.ProductPrice * float64(item.Quantity),
			}
			result, err := collection.InsertOne(ctx, newCart)
			if err != nil {
				logger.Error("Failed to create new cart", "error", err)
				return nil, err
			}

			if insertedId, ok := result.InsertedID.(bson.ObjectID); ok {
				newCart.Id = insertedId
			}
			logger.Info("New cart created successfully",
				"customer_id", customerID,
				"cart_id", newCart.Id.Hex(),
				"total_items", newCart.TotalItems,
				"total_price", newCart.TotalPrice,
			)
			return newCart, nil
		}
		logger.Error("Failed to find existing cart", "error", err)
		return nil, err
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
		logger.Error("Failed to update existing cart", "error", err)
		return nil, err
	}
	return &existingCart, nil

}

func (cr *cartRepositoryImpl) GetCartByCustomerId(ctx context.Context, customerID string) (*domain.Cart, error) {
	collection := cr.conn.Collection("carts")
	var cart domain.Cart
	err := collection.FindOne(ctx, bson.M{"customer_id": customerID}).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Error("Customer's cart is empty", "customer_id", customerID)
			return nil, err
		}
		logger.Error("Failed to find cart for customer", "customer_id", customerID, "error", err)
		return nil, err
	}
	return &cart, nil
}

func (cr *cartRepositoryImpl) UpdateCartItem(ctx context.Context, customerID string, item *domain.CartItem) (*domain.Cart, error) {
	collection := cr.conn.Collection("carts")
	var existingCart domain.Cart
	err := collection.FindOne(ctx, bson.M{"customer_id": customerID}).Decode(&existingCart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Error("Customer's cart is empty", "customer_id", customerID)
			return nil, err
		}
		logger.Error("Failed to find cart for customer", "customer_id", customerID, "error", err)
		return nil, err
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
		logger.Error("Item not found in cart", "product_id", item.ProductID, "customer_id", customerID)
		return nil, err
	}

	cr.recalCartTotals(&existingCart)
	_, err = collection.UpdateOne(ctx, bson.M{"_id": existingCart.Id}, bson.M{"$set": existingCart})
	if err != nil {
		logger.Error("Failed to update cart", "error", err)
		return nil, err
	}
	return &existingCart, nil
}

func (cr *cartRepositoryImpl) RemoveCartItem(ctx context.Context, customerID string, productID string) (*domain.Cart, error) {
	collection := cr.conn.Collection("carts")
	var existingCart domain.Cart
	err := collection.FindOne(ctx, bson.M{"customer_id": customerID}).Decode(&existingCart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Error("Customer's cart is empty", "customer_id", customerID)
			return nil, err
		}
		logger.Error("Failed to find cart for customer", "customer_id", customerID, "error", err)
		return nil, err
	}

	var updatedItems []*domain.CartItem
	for _, item := range existingCart.Items {
		if item.ProductID != productID {
			updatedItems = append(updatedItems, item)
		}
	}

	if len(updatedItems) == len(existingCart.Items) {
		logger.Error("Item not found in cart", "product_id", productID, "customer_id", customerID)
		return nil, err
	}

	existingCart.Items = updatedItems
	cr.recalCartTotals(&existingCart)

	_, err = collection.UpdateOne(ctx, bson.M{"_id": existingCart.Id}, bson.M{"$set": existingCart})
	if err != nil {
		logger.Error("Failed to update cart after removing item", "error", err)
		return nil, err
	}
	return &existingCart, nil
}

func (cr *cartRepositoryImpl) ClearCart(ctx context.Context, customerID string) error {
	collection := cr.conn.Collection("carts")
	result, err := collection.DeleteOne(ctx, bson.M{"customer_id": customerID})
	if err != nil {
		logger.Error("Failed to clear cart for customer", "customer_id", customerID, "error", err)
		return err
	}
	if result.DeletedCount == 0 {
		logger.Warn("Customer's cart is already empty", "customer_id", customerID)
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

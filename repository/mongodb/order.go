package mongodb

import (
	"context"
	"fmt"
	"intern-project-v2/domain"
	"intern-project-v2/logger"

	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var _ domain.OrderRepository = (*orderRepositoryImpl)(nil)

type orderRepositoryImpl struct {
	conn *mongo.Database
}

func NewOrderRepository(db *mongo.Database) domain.OrderRepository {
	return &orderRepositoryImpl{
		conn: db,
	}
}

func (or *orderRepositoryImpl) GetAll(ctx context.Context) ([]*domain.Order, error) {
	collection := or.conn.Collection("orders")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []*domain.Order
	for cursor.Next(ctx) {
		var order domain.Order
		if err := cursor.Decode(&order); err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (or *orderRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	collection := or.conn.Collection("orders")
	var order domain.Order
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		logger.Error("Invalid ID format", "id", id, "error", err)
		return nil, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Error("Order not found", "id", id)
			return nil, err
		}
		return nil, err
	}

	return &order, nil
}

func (or *orderRepositoryImpl) Create(ctx context.Context, order *domain.OrderRequest) (*domain.Order, error) {
	collection := or.conn.Collection("orders")
	newOrder := &domain.Order{
		CustomerId:  order.CustomerId,
		ProductIds:  order.ProductIds,
		TotalAmount: order.TotalAmount,
		CreatedAt:   time.Now(),
	}

	result, err := collection.InsertOne(ctx, newOrder)
	if err != nil {
		return nil, err
	}

	insertedID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		logger.Error("Failed to convert inserted ID to ObjectID", "insertedID", result.InsertedID)
		return nil, fmt.Errorf("failed to convert inserted ID to ObjectID: %v", result.InsertedID)
	}
	newOrder.Id = insertedID
	return newOrder, nil
}

func (or *orderRepositoryImpl) Update(ctx context.Context, id string, orderReq *domain.OrderRequest) (*domain.Order, error) {
	collection := or.conn.Collection("orders")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		logger.Error("Invalid ID format", "id", id, "error", err)
		return nil, err
	}

	updateFields := bson.M{}
	if orderReq.CustomerId != "" {
		updateFields["customer_id"] = orderReq.CustomerId
	}
	if len(orderReq.ProductIds) > 0 {
		updateFields["productids"] = orderReq.ProductIds
	}
	if orderReq.TotalAmount > 0 {
		updateFields["totalamount"] = orderReq.TotalAmount
	}

	update := bson.M{"$set": updateFields}
	otps := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update, otps)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			logger.Error("Order not found", "id", id)
			return nil, err
		}
		logger.Error("Failed to update order", "id", id, "error", result.Err())
		return nil, err
	}

	var updatedOrder domain.Order
	if err := result.Decode(&updatedOrder); err != nil {
		logger.Error("Failed to decode updated order", "id", id, "error", err)
		return nil, err
	}

	return &updatedOrder, nil
}
func (or *orderRepositoryImpl) Delete(ctx context.Context, id string) (*domain.Order, error) {
	collection := or.conn.Collection("orders")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		logger.Error("Invalid ID format", "id", id, "error", err)
		return nil, err
	}

	result := collection.FindOneAndDelete(ctx, bson.M{"_id": objectID})
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			logger.Error("Order not found for deletion", "id", id)
			return nil, err
		}
		logger.Error("Failed to delete order", "id", id, "error", result.Err())
		return nil, err
	}

	var deletedOrder domain.Order
	if err := result.Decode(&deletedOrder); err != nil {
		logger.Error("Failed to decode deleted order", "id", id, "error", err)
		return nil, err
	}

	return &deletedOrder, nil
}

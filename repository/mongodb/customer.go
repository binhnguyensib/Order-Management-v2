package mongodb

import (
	"context"
	"fmt"
	"intern-project-v2/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ domain.CustomerRepository = (*customerRepositoryImpl)(nil)

type customerRepositoryImpl struct {
	conn *mongo.Database
}

func NewCustomerRepository(db *mongo.Database) *customerRepositoryImpl {
	return &customerRepositoryImpl{
		conn: db,
	}
}

func (cr *customerRepositoryImpl) GetAll(ctx context.Context) ([]*domain.Customer, error) {
	collection := cr.conn.Collection("customers")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var customers []*domain.Customer
	for cursor.Next(context.TODO()) {
		var customer domain.Customer
		if err := cursor.Decode(&customer); err != nil {
			fmt.Printf("cusomter %v", customer)
			return nil, err
		}
		customers = append(customers, &customer)
	}

	return customers, nil
}

func (cr *customerRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	collection := cr.conn.Collection("customers")
	var customer domain.Customer
	ObjectID, ok := bson.ObjectIDFromHex(id)
	if ok != nil {
		return nil, fmt.Errorf("invalid ID format: %s", id)
	}
	err := collection.FindOne(ctx, bson.M{"_id": ObjectID}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("customer with ID %s not found", id)
		}
		return nil, err
	}
	return &customer, nil
}

func (cr *customerRepositoryImpl) Create(ctx context.Context, customer *domain.CustomerRequest) (*domain.Customer, error) {
	collection := cr.conn.Collection("customers")
	result, err := collection.InsertOne(ctx, customer)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %v", err)
	}
	customerID, ok := result.InsertedID.(bson.ObjectID)

	if !ok {
		return nil, fmt.Errorf("failed to convert inserted ID to ObjectID")
	}

	createdCustomer := &domain.Customer{
		Id:    customerID,
		Name:  customer.Name,
		Email: customer.Email,
		Phone: customer.Phone,
	}

	return createdCustomer, nil
}

func (cr *customerRepositoryImpl) Update(ctx context.Context, id string, customerReq *domain.CustomerRequest) (*domain.Customer, error) {
	collection := cr.conn.Collection("customers")
	ObjectID, ok := bson.ObjectIDFromHex(id)
	if ok != nil {
		return nil, fmt.Errorf("invalid ID format: %s", id)
	}

	updateFields := bson.M{}
	if customerReq.Name != "" {
		updateFields["name"] = customerReq.Name
	}
	if customerReq.Email != "" {
		updateFields["email"] = customerReq.Email
	}
	if customerReq.Phone != "" {
		updateFields["phone"] = customerReq.Phone
	}

	update := bson.M{"$set": updateFields}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": ObjectID}, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update customer: %v", err)
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("customer with ID %s not found", id)
	}

	return cr.GetByID(ctx, id)
}

func (cr *customerRepositoryImpl) Delete(ctx context.Context, id string) (*domain.Customer, error) {
	collection := cr.conn.Collection("customers")
	ObjectID, ok := bson.ObjectIDFromHex(id)
	if ok != nil {
		return nil, fmt.Errorf("invalid ID format: %s", id)
	}

	var deletedCustomer domain.Customer

	err := collection.FindOneAndDelete(ctx, bson.M{"_id": ObjectID}).Decode(&deletedCustomer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("customer with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to delete customer: %v", err)
	}
	return &deletedCustomer, nil
}

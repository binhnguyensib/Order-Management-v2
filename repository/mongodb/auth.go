package mongodb

import (
	"context"
	"intern-project-v2/domain"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ domain.AuthRepository = (*authRepositoryImpl)(nil)

type authRepositoryImpl struct {
	db *mongo.Database
}

func NewAuthRepository(db *mongo.Database) domain.AuthRepository {
	return &authRepositoryImpl{
		db: db,
	}
}
func (ar *authRepositoryImpl) Register(ctx context.Context, customer *domain.Customer) error {
	collection := ar.db.Collection("customers")
	_, err := collection.InsertOne(ctx, customer)
	if err != nil {
		return err
	}
	return nil
}
func (ar *authRepositoryImpl) Login(ctx context.Context, email string) (*domain.Customer, error) {
	collection := ar.db.Collection("customers")
	var customer domain.Customer
	err := collection.FindOne(ctx, map[string]interface{}{"email": email}).Decode(&customer)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

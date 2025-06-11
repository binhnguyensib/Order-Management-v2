package repository

import (
	"context"
	"intern-project-v2/domain"
)

type ProductRepository interface {
	GetAll(ctx context.Context) ([]*domain.Product, error)
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	Create(ctx context.Context, product *domain.ProductRequest) (*domain.Product, error)
	Update(ctx context.Context, id string, productReq *domain.ProductRequest) (*domain.Product, error)
	Delete(ctx context.Context, id string) (*domain.Product, error)
}

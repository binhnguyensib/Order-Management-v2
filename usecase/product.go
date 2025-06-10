package usecase

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

type ProductUsecase struct {
	productRepo ProductRepository
}

func NewProductUsecase(productRepo ProductRepository) *ProductUsecase {
	return &ProductUsecase{
		productRepo: productRepo,
	}
}

func (pu *ProductUsecase) GetAll(ctx context.Context) ([]*domain.Product, error) {
	products, err := pu.productRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil

}

func (pu *ProductUsecase) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	product, err := pu.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (pu *ProductUsecase) Create(ctx context.Context, product *domain.ProductRequest) (*domain.Product, error) {
	productCreated, err := pu.productRepo.Create(ctx, product)
	if err != nil {
		return nil, err
	}
	return productCreated, nil
}
func (pu *ProductUsecase) Update(ctx context.Context, id string, productReq *domain.ProductRequest) (*domain.Product, error) {
	productUpdated, err := pu.productRepo.Update(ctx, id, productReq)
	if err != nil {
		return nil, err
	}
	return productUpdated, nil
}

func (pu *ProductUsecase) Delete(ctx context.Context, id string) (*domain.Product, error) {
	productDeleted, err := pu.productRepo.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	return productDeleted, nil
}

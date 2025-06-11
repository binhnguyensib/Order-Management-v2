package usecase

import (
	"context"
	"intern-project-v2/domain"
)

type productUsecaseImpl struct {
	productRepo domain.ProductRepository
}

func NewProductUsecase(productRepo domain.ProductRepository) *productUsecaseImpl {
	return &productUsecaseImpl{
		productRepo: productRepo,
	}
}

func (pu *productUsecaseImpl) GetAll(ctx context.Context) ([]*domain.Product, error) {
	products, err := pu.productRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil

}

func (pu *productUsecaseImpl) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	product, err := pu.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (pu *productUsecaseImpl) Create(ctx context.Context, product *domain.ProductRequest) (*domain.Product, error) {
	productCreated, err := pu.productRepo.Create(ctx, product)
	if err != nil {
		return nil, err
	}
	return productCreated, nil
}
func (pu *productUsecaseImpl) Update(ctx context.Context, id string, productReq *domain.ProductRequest) (*domain.Product, error) {
	productUpdated, err := pu.productRepo.Update(ctx, id, productReq)
	if err != nil {
		return nil, err
	}
	return productUpdated, nil
}

func (pu *productUsecaseImpl) Delete(ctx context.Context, id string) (*domain.Product, error) {
	productDeleted, err := pu.productRepo.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	return productDeleted, nil
}

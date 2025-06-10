package usecase

import (
	"context"
	"intern-project-v2/domain"
)

type OrderRepository interface {
	GetAll(ctx context.Context) ([]*domain.Order, error)
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	Create(ctx context.Context, order *domain.OrderRequest) (*domain.Order, error)
	Update(ctx context.Context, id string, orderReq *domain.OrderRequest) (*domain.Order, error)
	Delete(ctx context.Context, id string) (*domain.Order, error)
}

type OrderUsecase struct {
	orderRepo OrderRepository
}

func NewOrderUsecase(orderRepo OrderRepository) *OrderUsecase {
	return &OrderUsecase{
		orderRepo: orderRepo,
	}
}
func (ou *OrderUsecase) GetAll(ctx context.Context) ([]*domain.Order, error) {
	orders, err := ou.orderRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
func (ou *OrderUsecase) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	order, err := ou.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return order, nil
}
func (ou *OrderUsecase) Create(ctx context.Context, order *domain.OrderRequest) (*domain.Order, error) {
	ord, err := ou.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, err
	}
	return ord, nil
}
func (ou *OrderUsecase) Update(ctx context.Context, id string, orderReq *domain.OrderRequest) (*domain.Order, error) {
	ord, err := ou.orderRepo.Update(ctx, id, orderReq)
	if err != nil {
		return nil, err
	}
	return ord, nil
}
func (ou *OrderUsecase) Delete(ctx context.Context, id string) (*domain.Order, error) {
	ord, err := ou.orderRepo.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	return ord, nil
}

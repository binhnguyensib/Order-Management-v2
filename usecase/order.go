package usecase

import (
	"context"
	"intern-project-v2/domain"
)

var _ domain.OrderUsecase = (*orderUsecaseImpl)(nil)

type orderUsecaseImpl struct {
	orderRepo domain.OrderRepository
}

func NewOrderUsecase(orderRepo domain.OrderRepository) *orderUsecaseImpl {
	return &orderUsecaseImpl{
		orderRepo: orderRepo,
	}
}
func (ou *orderUsecaseImpl) GetAll(ctx context.Context) ([]*domain.Order, error) {
	orders, err := ou.orderRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
func (ou *orderUsecaseImpl) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	order, err := ou.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return order, nil
}
func (ou *orderUsecaseImpl) Create(ctx context.Context, order *domain.OrderRequest) (*domain.Order, error) {
	ord, err := ou.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, err
	}
	return ord, nil
}
func (ou *orderUsecaseImpl) Update(ctx context.Context, id string, orderReq *domain.OrderRequest) (*domain.Order, error) {
	ord, err := ou.orderRepo.Update(ctx, id, orderReq)
	if err != nil {
		return nil, err
	}
	return ord, nil
}
func (ou *orderUsecaseImpl) Delete(ctx context.Context, id string) (*domain.Order, error) {
	ord, err := ou.orderRepo.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	return ord, nil
}

package usecase

import (
	"context"
	"intern-project-v2/domain"
)

type CustomerRepository interface {
	GetAll(ctx context.Context) ([]*domain.Customer, error)
	GetByID(ctx context.Context, id string) (*domain.Customer, error)
	Create(ctx context.Context, customer *domain.CustomerRequest) (*domain.Customer, error)
	Update(ctx context.Context, id string, customerReq *domain.CustomerRequest) (*domain.Customer, error)
	Delete(ctx context.Context, id string) (*domain.Customer, error)
}

type CustomerUsecase struct {
	customerRepo CustomerRepository
}

func NewCustomerUsecase(customerRepo CustomerRepository) *CustomerUsecase {
	return &CustomerUsecase{
		customerRepo: customerRepo,
	}
}
func (cu *CustomerUsecase) GetAll(ctx context.Context) ([]*domain.Customer, error) {
	customers, err := cu.customerRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return customers, nil
}

func (cu *CustomerUsecase) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	customer, err := cu.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (cu *CustomerUsecase) Create(ctx context.Context, customer *domain.CustomerRequest) (*domain.Customer, error) {
	cus, err := cu.customerRepo.Create(ctx, customer)
	if err != nil {
		return nil, err
	}
	return cus, nil
}

func (cu *CustomerUsecase) Update(ctx context.Context, id string, customerReq *domain.CustomerRequest) (*domain.Customer, error) {
	cus, err := cu.customerRepo.Update(ctx, id, customerReq)
	if err != nil {
		return nil, err
	}
	return cus, nil
}

func (cu *CustomerUsecase) Delete(ctx context.Context, id string) (*domain.Customer, error) {
	cus, err := cu.customerRepo.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	return cus, nil
}

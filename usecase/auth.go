package usecase

import (
	"context"
	"intern-project-v2/domain"
	"intern-project-v2/utils"
)

var _ domain.AuthUsecase = (*authUsecaseImpl)(nil)

type authUsecaseImpl struct {
	authRepo domain.AuthRepository
}

func NewAuthUsecase(authRepo domain.AuthRepository) domain.AuthUsecase {
	return &authUsecaseImpl{
		authRepo: authRepo,
	}
}
func (au *authUsecaseImpl) Register(ctx context.Context, customer *domain.CustomerRegiser) (*domain.Customer, error) {
	cust := &domain.Customer{
		Name:     customer.Name,
		Email:    customer.Email,
		Password: customer.Password,
		Phone:    customer.Phone,
	}

	if err := cust.HashPassword(); err != nil {
		return nil, err
	}

	if err := au.authRepo.Register(ctx, cust); err != nil {
		return nil, err
	}

	return cust, nil
}
func (au *authUsecaseImpl) Login(ctx context.Context, customer *domain.CustomerLogin) (*domain.Customer, string, error) {
	token, err := utils.GenerateJWT(customer.Email)
	if err != nil {
		return nil, "", err
	}
	cust, err := au.authRepo.Login(ctx, customer.Email)
	if err != nil {
		return nil, "", err
	}

	return cust, token, nil
}

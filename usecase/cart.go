package usecase

import (
	"context"
	"intern-project-v2/domain"
)

var _ domain.CartUsecase = (*cartUsecaseImpl)(nil)

type cartUsecaseImpl struct {
	cartRepo     domain.CartRepository
	productRepo  domain.ProductRepository
	customerRepo domain.CustomerRepository
}

func NewCartUsecase(
	cartRepo domain.CartRepository,
	productRepo domain.ProductRepository,
	customerRepo domain.CustomerRepository,
) domain.CartUsecase {
	return &cartUsecaseImpl{
		cartRepo:     cartRepo,
		productRepo:  productRepo,
		customerRepo: customerRepo,
	}
}

func (cu *cartUsecaseImpl) AddToCart(ctx context.Context, customerID string, cartItemReq *domain.CartItemRequest) (*domain.Cart, error) {
	productInfo, err := cu.productRepo.GetByID(ctx, cartItemReq.ProductID)
	if err != nil {
		return nil, err
	}

	cartItem := &domain.CartItem{
		ProductID:    cartItemReq.ProductID,
		ProductName:  cartItemReq.ProductName,
		Quantity:     cartItemReq.Quantity,
		ProductPrice: productInfo.Price,
		Subtotal:     float64(cartItemReq.Quantity) * productInfo.Price,
	}
	cart, err := cu.cartRepo.AddToCart(ctx, customerID, cartItem)
	if err != nil {
		return nil, err
	}
	return cart, nil
}

func (cu *cartUsecaseImpl) GetCartByCustomerId(ctx context.Context, customerID string) (*domain.Cart, error) {
	cart, err := cu.cartRepo.GetCartByCustomerId(ctx, customerID)
	if err != nil {
		return nil, err
	}
	return cart, nil
}

func (cu *cartUsecaseImpl) UpdateCartItem(ctx context.Context, customerID string, cartItemReq *domain.CartItemRequest) (*domain.Cart, error) {
	productInfo, err := cu.productRepo.GetByID(ctx, cartItemReq.ProductID)
	if err != nil {
		return nil, err
	}

	cartItem := &domain.CartItem{
		ProductID:    cartItemReq.ProductID,
		ProductName:  cartItemReq.ProductName,
		Quantity:     cartItemReq.Quantity,
		ProductPrice: productInfo.Price,
		Subtotal:     float64(cartItemReq.Quantity) * productInfo.Price,
	}
	cart, err := cu.cartRepo.UpdateCartItem(ctx, customerID, cartItem)
	if err != nil {
		return nil, err
	}
	return cart, nil
}

func (cu *cartUsecaseImpl) RemoveCartItem(ctx context.Context, customerID string, productID string) (*domain.Cart, error) {
	cart, err := cu.cartRepo.RemoveCartItem(ctx, customerID, productID)
	if err != nil {
		return nil, err
	}
	return cart, nil
}

func (cu *cartUsecaseImpl) ClearCart(ctx context.Context, customerID string) error {
	err := cu.cartRepo.ClearCart(ctx, customerID)
	if err != nil {
		return err
	}
	return nil
}

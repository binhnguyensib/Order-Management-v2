package domain

import (
	"context"
)

type ProductUsecase interface {
	GetAll(ctx context.Context) ([]*Product, error)
	GetByID(ctx context.Context, id string) (*Product, error)
	Create(ctx context.Context, product *ProductRequest) (*Product, error)
	Update(ctx context.Context, id string, productReq *ProductRequest) (*Product, error)
	Delete(ctx context.Context, id string) (*Product, error)
}

type ProductRepository interface {
	GetAll(ctx context.Context) ([]*Product, error)
	GetByID(ctx context.Context, id string) (*Product, error)
	Create(ctx context.Context, product *ProductRequest) (*Product, error)
	Update(ctx context.Context, id string, productReq *ProductRequest) (*Product, error)
	Delete(ctx context.Context, id string) (*Product, error)
}

type CustomerUsecase interface {
	GetAll(ctx context.Context) ([]*Customer, error)
	GetByID(ctx context.Context, id string) (*Customer, error)
	Create(ctx context.Context, customer *CustomerRequest) (*Customer, error)
	Update(ctx context.Context, id string, customerReq *CustomerRequest) (*Customer, error)
	Delete(ctx context.Context, id string) (*Customer, error)
}

type CustomerRepository interface {
	GetAll(ctx context.Context) ([]*Customer, error)
	GetByID(ctx context.Context, id string) (*Customer, error)
	Create(ctx context.Context, customer *CustomerRequest) (*Customer, error)
	Update(ctx context.Context, id string, customerReq *CustomerRequest) (*Customer, error)
	Delete(ctx context.Context, id string) (*Customer, error)
}

type OrderUsecase interface {
	GetAll(ctx context.Context) ([]*Order, error)
	GetByID(ctx context.Context, id string) (*Order, error)
	Create(ctx context.Context, order *OrderRequest) (*Order, error)
	Update(ctx context.Context, id string, orderReq *OrderRequest) (*Order, error)
	Delete(ctx context.Context, id string) (*Order, error)
}

type OrderRepository interface {
	GetAll(ctx context.Context) ([]*Order, error)
	GetByID(ctx context.Context, id string) (*Order, error)
	Create(ctx context.Context, order *OrderRequest) (*Order, error)
	Update(ctx context.Context, id string, orderReq *OrderRequest) (*Order, error)
	Delete(ctx context.Context, id string) (*Order, error)
}

type CartUsecase interface {
	AddToCart(ctx context.Context, customerID string, cart *CartItemRequest) (*Cart, error)
	GetCartByCustomerId(ctx context.Context, customerID string) (*Cart, error)
	UpdateCartItem(ctx context.Context, customerID string, cartItem *CartItemRequest) (*Cart, error)
	RemoveCartItem(ctx context.Context, customerID string, productID string) (*Cart, error)
	ClearCart(ctx context.Context, customerID string) error
}

type CartRepository interface {
	AddToCart(ctx context.Context, customerID string, cartItem *CartItem) (*Cart, error)
	GetCartByCustomerId(ctx context.Context, customerID string) (*Cart, error)
	UpdateCartItem(ctx context.Context, customerID string, cartItem *CartItem) (*Cart, error)
	RemoveCartItem(ctx context.Context, customerID string, productID string) (*Cart, error)
	ClearCart(ctx context.Context, customerID string) error
}

type AuthUsecase interface {
	Register(ctx context.Context, req *CustomerRegiser) (*Customer, error)
	Login(ctx context.Context, req *CustomerLogin) (*Customer, string, error)
}

type AuthRepository interface {
	Register(ctx context.Context, customer *Customer) error
	Login(ctx context.Context, email string) (*Customer, error)
}

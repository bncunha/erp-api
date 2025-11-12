package domain

import "context"

type GetProductsInput struct {
	SellerId *float64
}

type GetAllProductsOutput struct {
	Product  Product
	Quantity float64
}

type ProductRepository interface {
	Create(ctx context.Context, product Product) (int64, error)
	Edit(ctx context.Context, product Product, id int64) (int64, error)
	GetById(ctx context.Context, id int64) (Product, error)
	GetAll(ctx context.Context, input GetProductsInput) ([]GetAllProductsOutput, error)
	Inactivate(ctx context.Context, id int64) error
}

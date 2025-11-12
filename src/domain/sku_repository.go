package domain

import "context"

type GetSkusInput struct {
	SellerId *float64
}

type SkuRepository interface {
	Create(ctx context.Context, sku Sku, productId int64) (int64, error)
	CreateMany(ctx context.Context, skus []Sku, productId int64) ([]int64, error)
	GetByProductId(ctx context.Context, productId int64) ([]Sku, error)
	Update(ctx context.Context, sku Sku) error
	GetById(ctx context.Context, id int64) (Sku, error)
	GetByManyIds(ctx context.Context, ids []int64) ([]Sku, error)
	GetAll(ctx context.Context, input GetSkusInput) ([]Sku, error)
	Inactivate(ctx context.Context, id int64) error
}

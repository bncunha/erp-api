package domain

import "context"

type CategoryRepository interface {
	Create(ctx context.Context, category Category) (int64, error)
	GetById(ctx context.Context, id int64) (Category, error)
	GetByName(ctx context.Context, name string) (Category, error)
	Update(ctx context.Context, category Category) error
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context) ([]Category, error)
}

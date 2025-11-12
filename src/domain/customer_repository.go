package domain

import "context"

type CustomerRepository interface {
	GetById(ctx context.Context, id int64) (Customer, error)
	GetAll(ctx context.Context) ([]Customer, error)
	Create(ctx context.Context, customer Customer) (int64, error)
	Edit(ctx context.Context, customer Customer, id int64) (int64, error)
	Inactivate(ctx context.Context, id int64) error
}

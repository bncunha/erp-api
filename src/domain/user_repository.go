package domain

import "context"

type GetAllUserInput struct {
	Role *Role
}

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (User, error)
	Create(ctx context.Context, user User) (int64, error)
	Update(ctx context.Context, user User) error
	Inactivate(ctx context.Context, id int64) error
	GetAll(ctx context.Context, input GetAllUserInput) ([]User, error)
	GetById(ctx context.Context, id int64) (User, error)
}

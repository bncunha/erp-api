package domain

import (
	"context"
	"database/sql"
)

type GetAllUserInput struct {
	Role *Role
}

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (User, error)
	Create(ctx context.Context, user User) (int64, error)
	CreateWithTx(ctx context.Context, tx *sql.Tx, user User) (int64, error)
	Update(ctx context.Context, user User) error
	Inactivate(ctx context.Context, id int64) error
	GetAll(ctx context.Context, input GetAllUserInput) ([]User, error)
	GetById(ctx context.Context, id int64) (User, error)
	UpdatePassword(ctx context.Context, user User, newPassword string) error
	GetByEmail(ctx context.Context, email string) (User, error)
}

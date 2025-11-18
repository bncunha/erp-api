package domain

import "context"

type UserTokenRepository interface {
	Create(ctx context.Context, userToken UserToken) (int64, error)
	GetLastActiveByCodeHash(ctx context.Context, codeHash string) (UserToken, error)
	SetUsedToken(ctx context.Context, codeHash string) error
	GetById(ctx context.Context, id int64) (UserToken, error)
}

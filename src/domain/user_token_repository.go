package domain

import "context"

type UserTokenRepository interface {
	Create(ctx context.Context, userToken UserToken) (int64, error)
	GetLastActiveByUuid(ctx context.Context, uuid string) (UserToken, error)
	SetUsedToken(ctx context.Context, userToken UserToken) error
	GetById(ctx context.Context, id int64) (UserToken, error)
}

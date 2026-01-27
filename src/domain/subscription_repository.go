package domain

import (
	"context"
	"database/sql"
)

type SubscriptionRepository interface {
	GetByCompanyId(ctx context.Context, companyId int64) (Subscription, error)
	Create(ctx context.Context, subscription Subscription) (int64, error)
	CreateWithTx(ctx context.Context, tx *sql.Tx, subscription Subscription) (int64, error)
	UpdateWithTx(ctx context.Context, tx *sql.Tx, subscription Subscription) error
}

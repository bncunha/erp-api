package domain

import "context"

type NewsRepository interface {
	GetLatestVisible(ctx context.Context, tenantId int64, role Role) (News, error)
}

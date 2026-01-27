package domain

import "context"

type PlanRepository interface {
	GetByName(ctx context.Context, name string) (Plan, error)
	GetActiveByName(ctx context.Context, name string) (Plan, error)
}

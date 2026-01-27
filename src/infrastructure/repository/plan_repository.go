package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
)

type planRepository struct {
	db *sql.DB
}

func NewPlanRepository(db *sql.DB) domain.PlanRepository {
	return &planRepository{db}
}

func (r *planRepository) GetByName(ctx context.Context, name string) (domain.Plan, error) {
	var plan domain.Plan
	query := `SELECT id, name, price, status FROM plans WHERE name = $1`
	err := r.db.QueryRowContext(ctx, query, name).Scan(&plan.Id, &plan.Name, &plan.Price, &plan.Status)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return plan, domain.ErrPlanNotFound
		}
		return plan, err
	}
	return plan, nil
}

func (r *planRepository) GetActiveByName(ctx context.Context, name string) (domain.Plan, error) {
	var plan domain.Plan
	query := `SELECT id, name, price, status FROM plans WHERE name = $1 AND status = $2`
	err := r.db.QueryRowContext(ctx, query, name, domain.PlanStatusActive).Scan(&plan.Id, &plan.Name, &plan.Price, &plan.Status)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return plan, domain.ErrPlanNotFound
		}
		return plan, err
	}
	return plan, nil
}

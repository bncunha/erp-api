package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
)

type billingPaymentRepository struct {
	db *sql.DB
}

func NewBillingPaymentRepository(db *sql.DB) domain.BillingPaymentRepository {
	return &billingPaymentRepository{db}
}

func (r *billingPaymentRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, payment domain.BillingPayment) (int64, error) {
	var insertedID int64
	query := `INSERT INTO subscriptions_payments (company_id, subscription_id, plan_id, provider, status, amount, paid_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := tx.QueryRowContext(ctx, query,
		payment.CompanyId,
		payment.SubscriptionId,
		payment.PlanId,
		payment.Provider,
		payment.Status,
		payment.Amount,
		payment.PaidAt,
	).Scan(&insertedID)
	if err != nil {
		return insertedID, err
	}
	return insertedID, nil
}

func (r *billingPaymentRepository) GetByCompanyId(ctx context.Context, companyId int64) ([]domain.BillingPaymentHistory, error) {
	history := make([]domain.BillingPaymentHistory, 0)
	query := `SELECT p.id, p.plan_id, pl.name, p.provider, p.status, p.amount, p.paid_at, p.created_at
        FROM subscriptions_payments p
        JOIN plans pl ON pl.id = p.plan_id
        WHERE p.company_id = $1
        ORDER BY p.created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, companyId)
	if err != nil {
		return history, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.BillingPaymentHistory
		if err = rows.Scan(
			&item.Id,
			&item.PlanId,
			&item.PlanName,
			&item.Provider,
			&item.Status,
			&item.Amount,
			&item.PaidAt,
			&item.CreatedAt,
		); err != nil {
			return history, err
		}
		history = append(history, item)
	}
	return history, nil
}

func (r *billingPaymentRepository) GetLastPaymentDateByCompanyId(ctx context.Context, companyId int64) (*time.Time, error) {
	var paidAt sql.NullTime
	query := `SELECT COALESCE(MAX(paid_at), MAX(created_at))
        FROM subscriptions_payments
        WHERE company_id = $1`
	err := r.db.QueryRowContext(ctx, query, companyId).Scan(&paidAt)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return nil, nil
		}
		return nil, err
	}
	if !paidAt.Valid {
		return nil, nil
	}
	return &paidAt.Time, nil
}

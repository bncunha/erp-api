package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
)

type subscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) domain.SubscriptionRepository {
	return &subscriptionRepository{db}
}

func (r *subscriptionRepository) GetByCompanyId(ctx context.Context, companyId int64) (domain.Subscription, error) {
	var subscription domain.Subscription
	query := `SELECT s.id, s.company_id, s.plan_id, p.name, p.price, s.status, s.current_period_end,
        s.provider_name, s.provider_subscription_id, s.provider_customer_id
        FROM subscriptions s
        JOIN plans p ON p.id = s.plan_id
        WHERE s.company_id = $1`
	err := r.db.QueryRowContext(ctx, query, companyId).Scan(
		&subscription.Id,
		&subscription.CompanyId,
		&subscription.PlanId,
		&subscription.PlanName,
		&subscription.PlanPrice,
		&subscription.Status,
		&subscription.CurrentPeriodEnd,
		&subscription.ProviderName,
		&subscription.ProviderSubId,
		&subscription.ProviderCustId,
	)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return subscription, domain.ErrSubscriptionNotFound
		}
		return subscription, err
	}
	return subscription, nil
}

func (r *subscriptionRepository) Create(ctx context.Context, subscription domain.Subscription) (int64, error) {
	var insertedID int64
	query := `INSERT INTO subscriptions (company_id, plan_id, status, current_period_end, provider_name, provider_subscription_id, provider_customer_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := r.db.QueryRowContext(ctx, query,
		subscription.CompanyId,
		subscription.PlanId,
		subscription.Status,
		subscription.CurrentPeriodEnd,
		subscription.ProviderName,
		subscription.ProviderSubId,
		subscription.ProviderCustId,
	).Scan(&insertedID)
	if err != nil {
		return insertedID, err
	}
	return insertedID, nil
}

func (r *subscriptionRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, subscription domain.Subscription) (int64, error) {
	var insertedID int64
	query := `INSERT INTO subscriptions (company_id, plan_id, status, current_period_end, provider_name, provider_subscription_id, provider_customer_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := tx.QueryRowContext(ctx, query,
		subscription.CompanyId,
		subscription.PlanId,
		subscription.Status,
		subscription.CurrentPeriodEnd,
		subscription.ProviderName,
		subscription.ProviderSubId,
		subscription.ProviderCustId,
	).Scan(&insertedID)
	if err != nil {
		return insertedID, err
	}
	return insertedID, nil
}

func (r *subscriptionRepository) UpdateWithTx(ctx context.Context, tx *sql.Tx, subscription domain.Subscription) error {
	query := `UPDATE subscriptions
        SET plan_id = $1,
            status = $2,
            current_period_end = $3,
            updated_at = NOW()
        WHERE id = $4 AND company_id = $5`
	_, err := tx.ExecContext(ctx, query,
		subscription.PlanId,
		subscription.Status,
		subscription.CurrentPeriodEnd,
		subscription.Id,
		subscription.CompanyId,
	)
	return err
}

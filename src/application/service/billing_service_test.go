package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/domain"
)

func TestBillingServiceCreatesTrialWhenMissing(t *testing.T) {
	now := time.Date(2026, 1, 10, 8, 0, 0, 0, time.UTC)
	planRepo := &stubPlanRepository{
		activeByName: map[string]domain.Plan{
			domain.PlanNameTrial: {Id: 1, Name: domain.PlanNameTrial, Status: domain.PlanStatusActive},
		},
	}
	subRepo := &stubSubscriptionRepository{getErr: domain.ErrSubscriptionNotFound}
	service := &billingService{
		planRepository:         planRepo,
		subscriptionRepository: subRepo,
		now:                    func() time.Time { return now },
	}

	ctx := context.WithValue(context.Background(), constants.TENANT_KEY, int64(10))
	status, err := service.GetStatus(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.PlanName != domain.PlanNameTrial {
		t.Fatalf("expected trial plan, got %s", status.PlanName)
	}
	if !status.CanWrite {
		t.Fatalf("expected can_write true for trial")
	}
	expectedEnd := now.AddDate(0, 0, 15)
	if !subRepo.created.CurrentPeriodEnd.Equal(expectedEnd) {
		t.Fatalf("expected trial end %v, got %v", expectedEnd, subRepo.created.CurrentPeriodEnd)
	}
}

func TestBillingServiceExpiredSubscriptionBlocksWrite(t *testing.T) {
	now := time.Date(2026, 1, 10, 8, 0, 0, 0, time.UTC)
	subRepo := &stubSubscriptionRepository{
		subscription: domain.Subscription{
			Id:               1,
			CompanyId:        10,
			PlanId:           2,
			PlanName:         domain.PlanNameBasic,
			Status:           domain.SubscriptionStatusActive,
			CurrentPeriodEnd: now.Add(-time.Hour),
		},
	}
	service := &billingService{
		subscriptionRepository: subRepo,
		now:                    func() time.Time { return now },
	}

	ctx := context.WithValue(context.Background(), constants.TENANT_KEY, int64(10))
	status, err := service.GetStatus(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.CanWrite {
		t.Fatalf("expected can_write false")
	}
	if status.Reason != BillingReasonPaymentOverdue {
		t.Fatalf("expected reason %s, got %s", BillingReasonPaymentOverdue, status.Reason)
	}
}

type stubPlanRepository struct {
	activeByName map[string]domain.Plan
	byName       map[string]domain.Plan
	err          error
}

func (s *stubPlanRepository) GetByName(ctx context.Context, name string) (domain.Plan, error) {
	if s.err != nil {
		return domain.Plan{}, s.err
	}
	if plan, ok := s.byName[name]; ok {
		return plan, nil
	}
	return domain.Plan{}, domain.ErrPlanNotFound
}

func (s *stubPlanRepository) GetActiveByName(ctx context.Context, name string) (domain.Plan, error) {
	if s.err != nil {
		return domain.Plan{}, s.err
	}
	if plan, ok := s.activeByName[name]; ok {
		return plan, nil
	}
	return domain.Plan{}, domain.ErrPlanNotFound
}

type stubSubscriptionRepository struct {
	subscription domain.Subscription
	getErr       error
	created      domain.Subscription
	updated      domain.Subscription
	updateErr    error
}

func (s *stubSubscriptionRepository) GetByCompanyId(ctx context.Context, companyId int64) (domain.Subscription, error) {
	if s.getErr != nil {
		return domain.Subscription{}, s.getErr
	}
	return s.subscription, nil
}

func (s *stubSubscriptionRepository) Create(ctx context.Context, subscription domain.Subscription) (int64, error) {
	s.created = subscription
	return 1, nil
}

func (s *stubSubscriptionRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, subscription domain.Subscription) (int64, error) {
	return s.Create(ctx, subscription)
}

func (s *stubSubscriptionRepository) UpdateWithTx(ctx context.Context, tx *sql.Tx, subscription domain.Subscription) error {
	s.updated = subscription
	return s.updateErr
}

type stubBillingPaymentRepository struct {
	created domain.BillingPayment
	err     error
	history []domain.BillingPaymentHistory
	lastAt  *time.Time
}

func (s *stubBillingPaymentRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, payment domain.BillingPayment) (int64, error) {
	s.created = payment
	if s.err != nil {
		return 0, s.err
	}
	return 1, nil
}

func (s *stubBillingPaymentRepository) GetByCompanyId(ctx context.Context, companyId int64) ([]domain.BillingPaymentHistory, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.history, nil
}

func (s *stubBillingPaymentRepository) GetLastPaymentDateByCompanyId(ctx context.Context, companyId int64) (*time.Time, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.lastAt, nil
}

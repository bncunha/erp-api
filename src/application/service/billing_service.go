package service

import (
	"context"
	"strings"
	"time"

	"github.com/bncunha/erp-api/src/application/errors"
	helper "github.com/bncunha/erp-api/src/application/helpers"
	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/domain"
)

const (
	BillingReasonTrialExpired   = "TRIAL_EXPIRED"
	BillingReasonPaymentOverdue = "PAYMENT_OVERDUE"
)

type BillingService interface {
	GetStatus(ctx context.Context) (output.BillingStatusOutput, error)
	GetSummary(ctx context.Context) (output.BillingSummaryOutput, error)
	GetPayments(ctx context.Context) ([]output.BillingPaymentOutput, error)
}

type billingService struct {
	planRepository         domain.PlanRepository
	subscriptionRepository domain.SubscriptionRepository
	paymentRepository      domain.BillingPaymentRepository
	txManager              transactionManager
	now                    func() time.Time
}

func NewBillingService(planRepository domain.PlanRepository, subscriptionRepository domain.SubscriptionRepository, paymentRepository domain.BillingPaymentRepository, txManager transactionManager) BillingService {
	return &billingService{
		planRepository:         planRepository,
		subscriptionRepository: subscriptionRepository,
		paymentRepository:      paymentRepository,
		txManager:              txManager,
		now:                    time.Now,
	}
}

func (s *billingService) GetStatus(ctx context.Context) (output.BillingStatusOutput, error) {
	tenantId, err := helper.GetTenantId(ctx)
	if err != nil {
		return output.BillingStatusOutput{}, err
	}
	subscription, err := s.getOrCreateSubscription(ctx, tenantId)
	if err != nil {
		return output.BillingStatusOutput{}, err
	}
	return s.toStatus(subscription), nil
}

func (s *billingService) GetSummary(ctx context.Context) (output.BillingSummaryOutput, error) {
	tenantId, err := helper.GetTenantId(ctx)
	if err != nil {
		return output.BillingSummaryOutput{}, err
	}
	subscription, err := s.getOrCreateSubscription(ctx, tenantId)
	if err != nil {
		return output.BillingSummaryOutput{}, err
	}
	lastPaymentAt, err := s.paymentRepository.GetLastPaymentDateByCompanyId(ctx, tenantId)
	if err != nil {
		return output.BillingSummaryOutput{}, err
	}
	return output.BillingSummaryOutput{
		PlanName:      subscription.PlanName,
		PlanPrice:     subscription.PlanPrice,
		LastPaymentAt: lastPaymentAt,
		NextPaymentAt: subscription.CurrentPeriodEnd,
		Status:        string(subscription.Status),
	}, nil
}

func (s *billingService) GetPayments(ctx context.Context) ([]output.BillingPaymentOutput, error) {
	tenantId, err := helper.GetTenantId(ctx)
	if err != nil {
		return nil, err
	}
	history, err := s.paymentRepository.GetByCompanyId(ctx, tenantId)
	if err != nil {
		return nil, err
	}
	outputs := make([]output.BillingPaymentOutput, 0, len(history))
	for _, item := range history {
		outputs = append(outputs, output.BillingPaymentOutput{
			Id:        item.Id,
			PlanName:  item.PlanName,
			Provider:  string(item.Provider),
			Status:    string(item.Status),
			Amount:    item.Amount,
			PaidAt:    item.PaidAt,
			CreatedAt: item.CreatedAt,
		})
	}
	return outputs, nil
}

func (s *billingService) getOrCreateSubscription(ctx context.Context, tenantId int64) (domain.Subscription, error) {
	subscription, err := s.subscriptionRepository.GetByCompanyId(ctx, tenantId)
	if err == nil {
		return subscription, nil
	}
	if !errors.Is(err, domain.ErrSubscriptionNotFound) {
		return subscription, err
	}

	trialPlan, err := s.planRepository.GetActiveByName(ctx, domain.PlanNameTrial)
	if err != nil {
		return subscription, errors.New("Plano TRIAL não encontrado ou inativo")
	}

	trialEnd := s.now().AddDate(0, 0, 15)
	subscription = domain.Subscription{
		CompanyId:        tenantId,
		PlanId:           trialPlan.Id,
		PlanName:         trialPlan.Name,
		Status:           domain.SubscriptionStatusActive,
		CurrentPeriodEnd: trialEnd,
	}
	id, err := s.subscriptionRepository.Create(ctx, subscription)
	if err != nil {
		return subscription, err
	}
	subscription.Id = id
	return subscription, nil
}

func (s *billingService) toStatus(subscription domain.Subscription) output.BillingStatusOutput {
	now := s.now()
	expired := now.After(subscription.CurrentPeriodEnd)
	isActive := subscription.Status == domain.SubscriptionStatusActive
	canWrite := isActive && !expired

	reason := ""
	if !canWrite {
		if strings.EqualFold(subscription.PlanName, domain.PlanNameTrial) && expired {
			reason = BillingReasonTrialExpired
		} else {
			reason = BillingReasonPaymentOverdue
		}
	}

	return output.BillingStatusOutput{
		PlanName:         subscription.PlanName,
		CurrentPeriodEnd: subscription.CurrentPeriodEnd,
		CanWrite:         canWrite,
		Reason:           reason,
	}
}

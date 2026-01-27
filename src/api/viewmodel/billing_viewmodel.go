package viewmodel

import (
	"time"

	"github.com/bncunha/erp-api/src/application/service/output"
)

type BillingStatusViewModel struct {
	PlanName         string    `json:"plan_name"`
	CurrentPeriodEnd time.Time `json:"current_period_end"`
	CanWrite         bool      `json:"can_write"`
	Reason           string    `json:"reason,omitempty"`
}

func ToBillingStatusViewModel(status output.BillingStatusOutput) BillingStatusViewModel {
	return BillingStatusViewModel{
		PlanName:         status.PlanName,
		CurrentPeriodEnd: status.CurrentPeriodEnd,
		CanWrite:         status.CanWrite,
		Reason:           status.Reason,
	}
}

type BillingSummaryViewModel struct {
	PlanName      string     `json:"plan_name"`
	PlanPrice     float64    `json:"plan_price"`
	LastPaymentAt *time.Time `json:"last_payment_at"`
	NextPaymentAt time.Time  `json:"next_payment_at"`
	Status        string     `json:"status"`
}

func ToBillingSummaryViewModel(summary output.BillingSummaryOutput) BillingSummaryViewModel {
	return BillingSummaryViewModel{
		PlanName:      summary.PlanName,
		PlanPrice:     summary.PlanPrice,
		LastPaymentAt: summary.LastPaymentAt,
		NextPaymentAt: summary.NextPaymentAt,
		Status:        summary.Status,
	}
}

type BillingPaymentViewModel struct {
	Id        int64      `json:"id"`
	PlanName  string     `json:"plan_name"`
	Provider  string     `json:"provider"`
	Status    string     `json:"status"`
	Amount    float64    `json:"amount"`
	PaidAt    *time.Time `json:"paid_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func ToBillingPaymentViewModels(payments []output.BillingPaymentOutput) []BillingPaymentViewModel {
	viewModels := make([]BillingPaymentViewModel, 0, len(payments))
	for _, payment := range payments {
		viewModels = append(viewModels, BillingPaymentViewModel{
			Id:        payment.Id,
			PlanName:  payment.PlanName,
			Provider:  payment.Provider,
			Status:    payment.Status,
			Amount:    payment.Amount,
			PaidAt:    payment.PaidAt,
			CreatedAt: payment.CreatedAt,
		})
	}
	return viewModels
}

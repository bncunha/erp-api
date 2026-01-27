package output

import "time"

type BillingStatusOutput struct {
	PlanName         string
	CurrentPeriodEnd time.Time
	CanWrite         bool
	Reason           string
}

type BillingSummaryOutput struct {
	PlanName      string
	PlanPrice     float64
	LastPaymentAt *time.Time
	NextPaymentAt time.Time
	Status        string
}

type BillingPaymentOutput struct {
	Id        int64
	PlanName  string
	Provider  string
	Status    string
	Amount    float64
	PaidAt    *time.Time
	CreatedAt time.Time
}

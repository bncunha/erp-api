package domain

import (
	"time"
)

type BillingPaymentStatus string

const (
	BillingPaymentStatusPaid     BillingPaymentStatus = "PAID"
	BillingPaymentStatusPending  BillingPaymentStatus = "PENDING"
	BillingPaymentStatusFailed   BillingPaymentStatus = "FAILED"
	BillingPaymentStatusRefunded BillingPaymentStatus = "REFUNDED"
)

type BillingPaymentProvider string

const (
	BillingPaymentProviderManual BillingPaymentProvider = "MANUAL"
)

type BillingPayment struct {
	Id             int64
	CompanyId      int64
	SubscriptionId *int64
	PlanId         int64
	Provider       BillingPaymentProvider
	Status         BillingPaymentStatus
	Amount         float64
	PaidAt         *time.Time
	CreatedAt      time.Time
}

type BillingPaymentHistory struct {
	Id        int64
	PlanId    int64
	PlanName  string
	Provider  BillingPaymentProvider
	Status    BillingPaymentStatus
	Amount    float64
	PaidAt    *time.Time
	CreatedAt time.Time
}

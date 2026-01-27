package domain

import (
	"errors"
	"time"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive   SubscriptionStatus = "ACTIVE"
	SubscriptionStatusPastDue  SubscriptionStatus = "PAST_DUE"
	SubscriptionStatusCanceled SubscriptionStatus = "CANCELED"
)

var ErrSubscriptionNotFound = errors.New("Subscription not found")

type Subscription struct {
	Id               int64
	CompanyId        int64
	PlanId           int64
	PlanName         string
	PlanPrice        float64
	Status           SubscriptionStatus
	CurrentPeriodEnd time.Time
	ProviderName     *string
	ProviderSubId    *string
	ProviderCustId   *string
}

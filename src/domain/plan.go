package domain

import "errors"

type PlanStatus string

const (
	PlanStatusActive   PlanStatus = "ACTIVE"
	PlanStatusInactive PlanStatus = "INACTIVE"
)

const (
	PlanNameTrial   = "TRIAL"
	PlanNameBasic   = "BASIC"
	PlanNamePremium = "PREMIUM"
)

type Plan struct {
	Id     int64
	Name   string
	Price  float64
	Status PlanStatus
}

var ErrPlanNotFound = errors.New("Plan not found")

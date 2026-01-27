package response

import "time"

type BillingForbiddenResponse struct {
	Code             string    `json:"code"`
	Message          string    `json:"message"`
	Plan             string    `json:"plan"`
	CurrentPeriodEnd time.Time `json:"current_period_end"`
}

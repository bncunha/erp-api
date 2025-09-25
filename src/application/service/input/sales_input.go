package input

import (
	"time"

	"github.com/bncunha/erp-api/src/domain"
)

type GetSalesInput struct {
	InitialDate   *time.Time
	FinalDate     *time.Time
	UserId        []int64
	CustomerId    []int64
	PaymentStatus *domain.PaymentStatus
}

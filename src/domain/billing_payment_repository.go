package domain

import (
	"context"
	"database/sql"
	"time"
)

type BillingPaymentRepository interface {
	CreateWithTx(ctx context.Context, tx *sql.Tx, payment BillingPayment) (int64, error)
	GetByCompanyId(ctx context.Context, companyId int64) ([]BillingPaymentHistory, error)
	GetLastPaymentDateByCompanyId(ctx context.Context, companyId int64) (*time.Time, error)
}

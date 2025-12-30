package domain

import (
	"context"
	"database/sql"
)

type LegalAcceptanceRepository interface {
	CreateWithTx(ctx context.Context, tx *sql.Tx, acceptance LegalAcceptance) (int64, error)
}

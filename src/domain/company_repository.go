package domain

import (
    "context"
    "database/sql"
)

type CompanyRepository interface {
    CreateWithTx(ctx context.Context, tx *sql.Tx, company Company) (int64, error)
}

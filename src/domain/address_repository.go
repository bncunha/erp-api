package domain

import (
    "context"
    "database/sql"
)

type AddressRepository interface {
    CreateWithTx(ctx context.Context, tx *sql.Tx, address Address) (int64, error)
}

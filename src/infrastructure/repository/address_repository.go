package repository

import (
    "context"
    "database/sql"

    "github.com/bncunha/erp-api/src/domain"
)

type addressRepository struct {
    db *sql.DB
}

func NewAddressRepository(db *sql.DB) domain.AddressRepository {
    return &addressRepository{db}
}

func (r *addressRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, address domain.Address) (int64, error) {
    query := `INSERT INTO addresses (street, neighborhood, number, city, uf, cep, tenant_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
    var id int64
    err := tx.QueryRowContext(ctx, query, address.Street, address.Neighborhood, address.Number, address.City, address.UF, address.Cep, address.TenantId).Scan(&id)
    if err != nil {
        return id, err
    }
    return id, nil
}

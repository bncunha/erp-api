package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
)

type CustomerRepository interface {
	GetById(ctx context.Context, id int64) (domain.Customer, error)
}

type customerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepository{db}
}

func (r *customerRepository) GetById(ctx context.Context, id int64) (domain.Customer, error) {
	var customer domain.Customer
	tenantId := ctx.Value(constants.TENANT_KEY)

	query := `SELECT id, name, phone_number FROM customers WHERE id = $1 AND deleted_at IS NULL AND tenant_id = $2`
	err := r.db.QueryRowContext(ctx, query, id, tenantId).Scan(&customer.Id, &customer.Name, &customer.PhoneNumber)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return customer, errors.New("Cliente não encontrado")
		}
		return customer, err
	}
	return customer, nil
}

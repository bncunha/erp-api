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
    GetAll(ctx context.Context) ([]domain.Customer, error)
    Create(ctx context.Context, customer domain.Customer) (int64, error)
}

type customerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
    return &customerRepository{db}
}

func (r *customerRepository) Create(ctx context.Context, customer domain.Customer) (int64, error) {
    var insertedID int64
    tenantId := ctx.Value(constants.TENANT_KEY)
    query := `INSERT INTO customers (name, phone_number, tenant_id) VALUES ($1, $2, $3) RETURNING id`
    err := r.db.QueryRowContext(ctx, query, customer.Name, customer.PhoneNumber, tenantId).Scan(&insertedID)
    if err != nil {
        return insertedID, err
    }
    return insertedID, nil
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

func (r *customerRepository) GetAll(ctx context.Context) ([]domain.Customer, error) {
	var customers []domain.Customer
	tenantId := ctx.Value(constants.TENANT_KEY)

	query := `SELECT id, name, phone_number FROM customers WHERE deleted_at IS NULL AND tenant_id = $1`
	rows, err := r.db.QueryContext(ctx, query, tenantId)
	if err != nil {
		return customers, err
	}
	defer rows.Close()

	for rows.Next() {
		var customer domain.Customer
		err = rows.Scan(&customer.Id, &customer.Name, &customer.PhoneNumber)
		if err != nil {
			return customers, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

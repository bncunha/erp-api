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
	Edit(ctx context.Context, customer domain.Customer, id int64) (int64, error)
	Inactivate(ctx context.Context, id int64) error
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

func (r *customerRepository) Edit(ctx context.Context, customer domain.Customer, id int64) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var updatedID int64

	query := `UPDATE customers SET name = $1, phone_number = $2 WHERE id = $3 AND tenant_id = $4 AND deleted_at IS NULL RETURNING id`
	err := r.db.QueryRowContext(ctx, query, customer.Name, customer.PhoneNumber, id, tenantId).Scan(&updatedID)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return updatedID, errors.New("Cliente não encontrado")
		}
		return updatedID, err
	}
	return updatedID, nil
}

func (r *customerRepository) Inactivate(ctx context.Context, id int64) error {
	tenantId := ctx.Value(constants.TENANT_KEY)
	query := `DELETE FROM customers WHERE id = $1 AND tenant_id = $2`
	result, err := r.db.ExecContext(ctx, query, id, tenantId)
	if err != nil {
		if errors.IsForeignKeyViolation(err) {
			return errors.New("Não é possível deletar o cliente pois existem registros associados.")
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Cliente não encontrado")
	}

	return nil
}

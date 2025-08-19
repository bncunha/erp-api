package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (domain.User, error)
	Create(ctx context.Context, user domain.User) (int64, error)
	Update(ctx context.Context, user domain.User) error
	Inactivate(ctx context.Context, id int64) error
	GetAll(ctx context.Context) ([]domain.User, error)
	GetById(ctx context.Context, id int64) (domain.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	var user domain.User
	tenantId := ctx.Value(constants.TENANT_KEY)

	query := `SELECT id, username, name, phone_number, password, role, tenant_id FROM users WHERE username = $1 AND tenant_id = $2 AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, username, tenantId).Scan(
		&user.Id,
		&user.Username,
		&user.Name,
		&user.PhoneNumber,
		&user.Password,
		&user.Role,
		&user.TenantId,
	)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return user, errors.New("Usuário não encontrado")
		}
		return user, err
	}
	return user, nil
}

func (r *userRepository) Create(ctx context.Context, user domain.User) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	query := `INSERT INTO users (username, name, phone_number, password, role, tenant_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var id int64
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Name, user.PhoneNumber, user.Password, user.Role, tenantId).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (r *userRepository) Update(ctx context.Context, user domain.User) error {
	tenantId := ctx.Value(constants.TENANT_KEY)
	query := `UPDATE users SET name = $1, phone_number = $2, password = $3, role = $4 WHERE id = $5 AND tenant_id = $6 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, user.Name, user.PhoneNumber, user.Password, user.Role, user.Id, tenantId)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Inactivate(ctx context.Context, id int64) error {
	tenantId := ctx.Value(constants.TENANT_KEY)
	query := `UPDATE users SET deleted_at = false WHERE id = $1 and tenant_id = $2`
	_, err := r.db.ExecContext(ctx, query, id, tenantId)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var users []domain.User

	query := `SELECT id, username, name, phone_number, password, role, tenant_id FROM users WHERE tenant_id = $1 AND deleted_at IS NULL ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query, tenantId)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.Id, &user.Username, &user.Name, &user.PhoneNumber, &user.Password, &user.Role, &user.TenantId)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, err
}

func (r *userRepository) GetById(ctx context.Context, id int64) (domain.User, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var user domain.User

	query := `SELECT id, username, name, phone_number, password, role, tenant_id FROM users WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, id, tenantId).Scan(&user.Id, &user.Username, &user.Name, &user.PhoneNumber, &user.Password, &user.Role, &user.TenantId)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return user, errors.New("User não encontrada")
		}
		return user, err
	}
	return user, nil
}

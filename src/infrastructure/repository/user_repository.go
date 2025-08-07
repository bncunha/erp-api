package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (domain.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	var user domain.User

	query := `SELECT id, username, name, phone_number, password, role, tenant_id FROM users WHERE username = $1`
	err := r.db.QueryRowContext(ctx, query, username).Scan(
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
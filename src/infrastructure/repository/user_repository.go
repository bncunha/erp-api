package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	var user domain.User
	query := `SELECT id, username, name, phone_number, password, role, tenant_id FROM users WHERE username = $1 AND deleted_at IS NULL`
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

func (r *userRepository) Create(ctx context.Context, user domain.User) (int64, error) {
	return r.create(ctx, r.db, nil, user)
}

func (r *userRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, user domain.User) (int64, error) {
	return r.create(ctx, nil, tx, user)
}

func (r *userRepository) create(ctx context.Context, db *sql.DB, tx *sql.Tx, user domain.User) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	query := `INSERT INTO users (username, name, phone_number, role, tenant_id, email, password) VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, '')) RETURNING id`
	var id int64
	executor := any(db)
	if tx != nil {
		executor = tx
	}

	switch exec := executor.(type) {
	case interface {
		QueryRowContext(context.Context, string, ...any) *sql.Row
	}:
		err := exec.QueryRowContext(ctx, query, user.Username, user.Name, user.PhoneNumber, user.Role, tenantId, user.Email, user.Password).Scan(&id)
		if err != nil {
			return id, err
		}
	default:
		return id, errors.New("executor not provided")
	}
	return id, nil
}

func (r *userRepository) Update(ctx context.Context, user domain.User) error {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var err error
	query := `UPDATE users SET name = $1, phone_number = $2, role = $3, username = $4, email = $7 WHERE id = $5 AND tenant_id = $6`

	_, err = r.db.ExecContext(ctx, query, user.Name, user.PhoneNumber, user.Role, user.Username, user.Id, tenantId, user.Email)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Inactivate(ctx context.Context, id int64) error {
	tenantId := ctx.Value(constants.TENANT_KEY)
	query := `DELETE FROM users WHERE id = $1 AND tenant_id = $2`
	result, err := r.db.ExecContext(ctx, query, id, tenantId)
	if err != nil {
		if errors.IsForeignKeyViolation(err) {
			return errors.New("Não é possível deletar o usuário pois existem registros associados.")
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Usuário não encontrado")
	}

	return nil
}

func (r *userRepository) GetAll(ctx context.Context, input domain.GetAllUserInput) ([]domain.User, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var users []domain.User

	query := `SELECT 
		id,
		username, 
		name, 
		phone_number, 
		role, 
		tenant_id,
		email
	FROM users 
	WHERE tenant_id = $1 AND deleted_at IS NULL AND ($2::text IS NULL OR role = $2)
	ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query, tenantId, input.Role)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.Id, &user.Username, &user.Name, &user.PhoneNumber, &user.Role, &user.TenantId, &user.Email)
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

	query := `SELECT id, username, name, phone_number, role, tenant_id, email FROM users WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, id, tenantId).Scan(&user.Id, &user.Username, &user.Name, &user.PhoneNumber, &user.Role, &user.TenantId, &user.Email)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return user, errors.New("Usuário não encontrado")
		}
		return user, err
	}
	return user, nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, user domain.User, newPassword string) error {
	query := `UPDATE users SET password = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, newPassword, user.Id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Usuário não encontrado")
	}

	return nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	query := `SELECT id, username, name, phone_number, role, tenant_id, email FROM users WHERE email = $1 AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.Username,
		&user.Name,
		&user.PhoneNumber,
		&user.Role,
		&user.TenantId,
		&user.Email,
	)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return user, errors.New("Usuário não encontrado")
		}
		return user, err
	}
	return user, nil
}

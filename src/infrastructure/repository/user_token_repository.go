package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
)

type userTokenRepository struct {
	db *sql.DB
}

func NewUserTokenRepository(db *sql.DB) domain.UserTokenRepository {
	return &userTokenRepository{db: db}
}

func (r *userTokenRepository) Create(ctx context.Context, userToken domain.UserToken) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)

	query := `INSERT INTO user_tokens (user_id, tenant_id, type, code_hash, expires_at, used_at, created_by, uuid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		userToken.User.Id,
		tenantId,
		userToken.Type,
		userToken.CodeHash,
		userToken.ExpiresAt,
		userToken.UsedAt,
		userToken.CreatedBy.Id,
		userToken.Uuid,
	).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (r *userTokenRepository) GetLastActiveByUuid(ctx context.Context, uuid string) (domain.UserToken, error) {
	var userToken domain.UserToken
	var userID int64
	usedAt := sql.NullTime{}
	createdBy := sql.NullInt64{}

	query := `
		SELECT ut.id, ut.user_id, ut.type, ut.code_hash, ut.expires_at, ut.used_at, ut.created_by,
			u.username, u.name, u.phone_number, u.role, u.tenant_id, u.email, ut.uuid
		FROM user_tokens ut
		INNER JOIN users u ON u.id = ut.user_id
		WHERE ut.uuid = $1
			AND ut.used_at IS NULL
			AND ut.expires_at >= NOW()
		ORDER BY ut.created_at DESC
		LIMIT 1`

	err := r.db.QueryRowContext(ctx, query, uuid).Scan(
		&userToken.Id,
		&userID,
		&userToken.Type,
		&userToken.CodeHash,
		&userToken.ExpiresAt,
		&usedAt,
		&createdBy,
		&userToken.User.Username,
		&userToken.User.Name,
		&userToken.User.PhoneNumber,
		&userToken.User.Role,
		&userToken.User.TenantId,
		&userToken.User.Email,
		&userToken.Uuid,
	)
	if err != nil {
		return userToken, err
	}

	userToken.User.Id = userID
	if usedAt.Valid {
		userToken.UsedAt = &usedAt.Time
	}
	if createdBy.Valid {
		userToken.CreatedBy.Id = createdBy.Int64
	}
	return userToken, nil
}

func (r *userTokenRepository) SetUsedToken(ctx context.Context, userToken domain.UserToken) error {
	query := `UPDATE user_tokens SET used_at = $2 WHERE id = $1 AND user_id = $3 AND used_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, userToken.Id, userToken.UsedAt, userToken.User.Id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Token não encontrado ou já utilizado")
	}

	return nil
}

func (r *userTokenRepository) GetById(ctx context.Context, id int64) (domain.UserToken, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var userToken domain.UserToken
	var userID int64
	usedAt := sql.NullTime{}
	createdBy := sql.NullInt64{}

	query := `SELECT id, user_id, type, code_hash, expires_at, used_at, created_by, uuid
		FROM user_tokens
		WHERE id = $1 AND ($2::bigint IS NULL OR tenant_id = $2::bigint)`

	err := r.db.QueryRowContext(ctx, query, id, tenantId).Scan(
		&userToken.Id,
		&userID,
		&userToken.Type,
		&userToken.CodeHash,
		&userToken.ExpiresAt,
		&usedAt,
		&createdBy,
		&userToken.Uuid,
	)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return userToken, errors.New("Token não encontrado")
		}
		return userToken, err
	}

	userToken.User.Id = userID
	if usedAt.Valid {
		userToken.UsedAt = &usedAt.Time
	}
	if createdBy.Valid {
		userToken.CreatedBy.Id = createdBy.Int64
	}
	return userToken, nil
}

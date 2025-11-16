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
	usedAt := sql.NullTime{}
	if !userToken.UsedAt.IsZero() {
		usedAt = sql.NullTime{Time: userToken.UsedAt, Valid: true}
	}

	createdBy := sql.NullInt64{}
	if userToken.CreatedBy.Id != 0 {
		createdBy = sql.NullInt64{Int64: userToken.CreatedBy.Id, Valid: true}
	}

	query := `INSERT INTO user_tokens (user_id, tenant_id, type, code_hash, expires_at, used_at, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		userToken.UserId.Id,
		tenantId,
		userToken.Type,
		userToken.CodeHash,
		userToken.ExpiresAt,
		usedAt,
		createdBy,
	).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (r *userTokenRepository) GetLastActiveByCodeHash(ctx context.Context, codeHash string) (domain.UserToken, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var userToken domain.UserToken
	var userID int64
	usedAt := sql.NullTime{}
	createdBy := sql.NullInt64{}

	query := `SELECT id, user_id, type, code_hash, expires_at, used_at, created_by, created_at
		FROM user_tokens
		WHERE code_hash = $1
			AND tenant_id = $2
			AND used_at IS NULL
			AND expires_at >= NOW()
		ORDER BY created_at DESC
		LIMIT 1`

	err := r.db.QueryRowContext(ctx, query, codeHash, tenantId).Scan(
		&userToken.Id,
		&userID,
		&userToken.Type,
		&userToken.CodeHash,
		&userToken.ExpiresAt,
		&usedAt,
		&createdBy,
		&userToken.CreatedAt,
	)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return userToken, errors.New("Token inválido ou expirado")
		}
		return userToken, err
	}

	userToken.UserId.Id = userID
	if usedAt.Valid {
		userToken.UsedAt = usedAt.Time
	}
	if createdBy.Valid {
		userToken.CreatedBy.Id = createdBy.Int64
	}
	return userToken, nil
}

func (r *userTokenRepository) SetUsedToken(ctx context.Context, codeHash string) error {
	tenantId := ctx.Value(constants.TENANT_KEY)
	query := `UPDATE user_tokens SET used_at = NOW() WHERE code_hash = $1 AND used_at IS NULL AND tenant_id = $2`
	result, err := r.db.ExecContext(ctx, query, codeHash, tenantId)
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

package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
)

type legalAcceptanceRepository struct {
	db *sql.DB
}

func NewLegalAcceptanceRepository(db *sql.DB) domain.LegalAcceptanceRepository {
	return &legalAcceptanceRepository{db: db}
}

func (r *legalAcceptanceRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, acceptance domain.LegalAcceptance) (int64, error) {
	if tx == nil {
		return 0, errors.New("transaction not provided")
	}

	query := `INSERT INTO legal_acceptances (user_id, tenant_id, legal_document_id, accepted)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
	var id int64
	err := tx.QueryRowContext(ctx, query, acceptance.UserId, acceptance.TenantId, acceptance.LegalDocumentId, acceptance.Accepted).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

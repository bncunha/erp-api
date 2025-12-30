package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
)

type legalDocumentRepository struct {
	db *sql.DB
}

func NewLegalDocumentRepository(db *sql.DB) domain.LegalDocumentRepository {
	return &legalDocumentRepository{db: db}
}

func (r *legalDocumentRepository) GetLastActiveByType(ctx context.Context, docType domain.LegalDocumentType) (domain.LegalDocument, error) {
	var document domain.LegalDocument
	query := `SELECT id, doc_type, doc_version, published_at, content_sha256, is_active
		FROM legal_documents
		WHERE doc_type = $1 AND is_active = true
		ORDER BY published_at DESC
		LIMIT 1`

	err := r.db.QueryRowContext(ctx, query, docType).Scan(
		&document.Id,
		&document.DocType,
		&document.DocVersion,
		&document.PublishedAt,
		&document.ContentSha256,
		&document.IsActive,
	)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return document, errors.New("Documento legal ativo nao encontrado")
		}
		return document, err
	}

	return document, nil
}

func (r *legalDocumentRepository) GetByTypeAndVersion(ctx context.Context, docType domain.LegalDocumentType, version string) (domain.LegalDocument, error) {
	var document domain.LegalDocument
	query := `SELECT id, doc_type, doc_version, published_at, content_sha256, is_active
		FROM legal_documents
		WHERE doc_type = $1 AND doc_version = $2
		LIMIT 1`

	err := r.db.QueryRowContext(ctx, query, docType, version).Scan(
		&document.Id,
		&document.DocType,
		&document.DocVersion,
		&document.PublishedAt,
		&document.ContentSha256,
		&document.IsActive,
	)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return document, errors.New("Documento legal nao encontrado")
		}
		return document, err
	}

	return document, nil
}

func (r *legalDocumentRepository) GetActiveByUser(ctx context.Context, userId int64) ([]domain.LegalTermStatus, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	terms := make([]domain.LegalTermStatus, 0)

	query := `SELECT d.doc_type, d.doc_version, COALESCE(a.accepted, false) AS accepted
		FROM (
			SELECT DISTINCT ON (doc_type) id, doc_type, doc_version
			FROM legal_documents
			WHERE is_active = true
			ORDER BY doc_type, published_at DESC
		) d
		LEFT JOIN legal_acceptances a
			ON a.legal_document_id = d.id
			AND a.user_id = $1
			AND a.tenant_id = $2
		ORDER BY d.doc_type ASC`

	rows, err := r.db.QueryContext(ctx, query, userId, tenantId)
	if err != nil {
		return terms, err
	}
	defer rows.Close()

	for rows.Next() {
		var term domain.LegalTermStatus
		if err = rows.Scan(&term.DocType, &term.DocVersion, &term.Accepted); err != nil {
			return terms, err
		}
		terms = append(terms, term)
	}

	return terms, err
}

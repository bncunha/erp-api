package domain

import "context"

type LegalDocumentRepository interface {
	GetLastActiveByType(ctx context.Context, docType LegalDocumentType) (LegalDocument, error)
	GetActiveByUser(ctx context.Context, userId int64) ([]LegalTermStatus, error)
}

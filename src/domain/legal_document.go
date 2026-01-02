package domain

import "time"

type LegalDocumentType string

const (
	LegalDocumentTypeTerms   LegalDocumentType = "TERMS"
	LegalDocumentTypePrivacy LegalDocumentType = "PRIVACY"
)

type LegalDocument struct {
	Id            int64
	DocType       LegalDocumentType
	DocVersion    string
	PublishedAt   time.Time
	ContentSha256 string
	IsActive      bool
}

type LegalAcceptance struct {
	Id              int64
	UserId          int64
	TenantId        int64
	LegalDocumentId int64
	Accepted        bool
}

type LegalTermStatus struct {
	DocType    LegalDocumentType
	DocVersion string
	Accepted   bool
}

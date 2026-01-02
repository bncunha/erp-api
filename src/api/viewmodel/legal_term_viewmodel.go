package viewmodel

import "github.com/bncunha/erp-api/src/domain"

type LegalTermViewModel struct {
	DocType    string `json:"doc_type"`
	DocVersion string `json:"doc_version"`
	Accepted   bool   `json:"accepted"`
}

func ToLegalTermViewModel(term domain.LegalTermStatus) LegalTermViewModel {
	return LegalTermViewModel{
		DocType:    string(term.DocType),
		DocVersion: term.DocVersion,
		Accepted:   term.Accepted,
	}
}

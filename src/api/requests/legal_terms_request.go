package request

import (
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/application/validator"
)

type AcceptLegalTermRequest struct {
	DocType    string `json:"doc_type" validate:"required,oneof=TERMS PRIVACY"`
	DocVersion string `json:"doc_version" validate:"required"`
	Accepted   bool   `json:"accepted" validate:"required"`
}

func (r *AcceptLegalTermRequest) Validate() error {
	if err := validator.Validate(r); err != nil {
		return err
	}
	if !r.Accepted {
		return errors.New("Aceite o termo para continuar.")
	}
	return nil
}

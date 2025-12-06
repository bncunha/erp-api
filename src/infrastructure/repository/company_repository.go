package repository

import (
    "context"
    "database/sql"

    "github.com/bncunha/erp-api/src/domain"
)

type companyRepository struct {
    db *sql.DB
}

func NewCompanyRepository(db *sql.DB) domain.CompanyRepository {
    return &companyRepository{db}
}

func (r *companyRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, company domain.Company) (int64, error) {
    query := `INSERT INTO companies (name, legal_name, cnpj, cpf, cellphone) VALUES ($1, $2, $3, $4, $5) RETURNING id`
    var id int64

    cnpj := sql.NullString{String: company.Cnpj, Valid: company.Cnpj != ""}
    cpf := sql.NullString{String: company.Cpf, Valid: company.Cpf != ""}
    cellphone := sql.NullString{String: company.Cellphone, Valid: company.Cellphone != ""}
    err := tx.QueryRowContext(ctx, query, company.Name, company.LegalName, cnpj, cpf, cellphone).Scan(&id)
    if err != nil {
        return id, err
    }
    return id, nil
}

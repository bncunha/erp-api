package repository

import (
	"context"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/bncunha/erp-api/src/application/constants"
)

func TestQuoteRepositoryGetCompanyProfileReturnsAddress(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewQuoteRepository(db)
	queryRegex := regexp.QuoteMeta(`SELECT c.name, c.legal_name, COALESCE(c.cnpj,''), COALESCE(c.email,''), COALESCE(c.cellphone,''), c.logo_url,
	COALESCE(a.street,'') || CASE WHEN a.number IS NOT NULL AND a.number <> '' THEN ', ' || a.number ELSE '' END ||
	CASE WHEN a.neighborhood IS NOT NULL AND a.neighborhood <> '' THEN ' - ' || a.neighborhood ELSE '' END ||
	CASE WHEN a.city IS NOT NULL AND a.city <> '' THEN ' - ' || a.city ELSE '' END ||
	CASE WHEN a.uf IS NOT NULL AND a.uf <> '' THEN '/' || a.uf ELSE '' END ||
	CASE WHEN a.cep IS NOT NULL AND a.cep <> '' THEN ' - CEP ' || a.cep ELSE '' END AS full_address
	FROM companies c
	LEFT JOIN LATERAL (
		SELECT street, number, neighborhood, city, uf, cep
		FROM addresses
		WHERE tenant_id=c.id
		ORDER BY id DESC
		LIMIT 1
	) a ON TRUE
	WHERE c.id=$1`)

	rows := sqlmock.NewRows([]string{"name", "legal_name", "cnpj", "email", "cellphone", "logo_url", "full_address"}).
		AddRow("ACME", "ACME LTDA", "123", "contato@acme.com", "11999999999", "logo.png", "Rua A, 10 - Centro - Sao Paulo/SP - CEP 01001000")

	mock.ExpectQuery(queryRegex).
		WithArgs(int64(2)).
		WillReturnRows(rows)

	ctx := context.WithValue(context.Background(), constants.TENANT_KEY, int64(2))
	profile, err := repo.GetCompanyProfile(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if profile.Address != "Rua A, 10 - Centro - Sao Paulo/SP - CEP 01001000" {
		t.Fatalf("unexpected address: %q", profile.Address)
	}
	if profile.LogoURL == nil || *profile.LogoURL != "logo.png" {
		t.Fatalf("unexpected logo_url: %+v", profile.LogoURL)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestQuoteRepositoryGetCompanyProfileNoLogoURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewQuoteRepository(db)
	rows := sqlmock.NewRows([]string{"name", "legal_name", "cnpj", "email", "cellphone", "logo_url", "full_address"}).
		AddRow("ACME", "ACME LTDA", "123", "contato@acme.com", "11999999999", nil, "Av B, 200 - Bairro - Campinas/SP - CEP 13000000")

	mock.ExpectQuery("SELECT c.name, c.legal_name").
		WithArgs(int64(7)).
		WillReturnRows(rows)

	ctx := context.WithValue(context.Background(), constants.TENANT_KEY, int64(7))
	profile, err := repo.GetCompanyProfile(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if profile.Address != "Av B, 200 - Bairro - Campinas/SP - CEP 13000000" {
		t.Fatalf("unexpected address: %q", profile.Address)
	}
	if profile.LogoURL != nil {
		t.Fatalf("expected nil logo_url, got: %+v", profile.LogoURL)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestQuoteRepositoryGetCompanyProfileWithoutAddressReturnsEmptyString(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewQuoteRepository(db)
	rows := sqlmock.NewRows([]string{"name", "legal_name", "cnpj", "email", "cellphone", "logo_url", "full_address"}).
		AddRow("ACME", "ACME LTDA", "123", "contato@acme.com", "11999999999", nil, "")

	mock.ExpectQuery("SELECT c.name, c.legal_name").
		WithArgs(int64(9)).
		WillReturnRows(rows)

	ctx := context.WithValue(context.Background(), constants.TENANT_KEY, int64(9))
	profile, err := repo.GetCompanyProfile(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if profile.Address != "" {
		t.Fatalf("expected empty address, got: %q", profile.Address)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

package helper

import (
	"context"
	"testing"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/domain"
)

func TestParseInt64(t *testing.T) {
	if got := ParseInt64("42"); got != 42 {
		t.Fatalf("expected 42, got %d", got)
	}
	if got := ParseInt64("invalid"); got != 0 {
		t.Fatalf("expected 0 for invalid input, got %d", got)
	}
}

func TestGetRole(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.ROLE_KEY, string(domain.UserRoleAdmin))
	if role := GetRole(ctx); role != domain.UserRoleAdmin {
		t.Fatalf("expected admin role, got %s", role)
	}
}

func TestGetTenantId(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.TENANT_KEY, int64(42))
	tenantID, err := GetTenantId(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tenantID != 42 {
		t.Fatalf("expected tenant id 42, got %d", tenantID)
	}

	ctx = context.WithValue(context.Background(), constants.TENANT_KEY, float64(55))
	tenantID, err = GetTenantId(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tenantID != 55 {
		t.Fatalf("expected tenant id 55, got %d", tenantID)
	}

	ctx = context.WithValue(context.Background(), constants.TENANT_KEY, "invalid")
	_, err = GetTenantId(ctx)
	if err == nil {
		t.Fatalf("expected error for invalid tenant id")
	}
}

func TestParseFloat(t *testing.T) {
	value, err := ParseFloat("12.50")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != 12.50 {
		t.Fatalf("expected 12.50, got %v", value)
	}

	if _, err := ParseFloat("invalid"); err == nil {
		t.Fatalf("expected error for invalid float")
	}
}

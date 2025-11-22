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

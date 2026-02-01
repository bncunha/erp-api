package helper

import "testing"

func TestGenerateAndParseJWT(t *testing.T) {
	token, err := GenerateJWT("user", 123, "ADMIN", 456)
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	username, tenant, role, userID, err := ParseJWT(token)
	if err != nil {
		t.Fatalf("unexpected error parsing token: %v", err)
	}

	if username != "user" {
		t.Fatalf("expected username 'user', got %s", username)
	}
	if int64(tenant) != 123 {
		t.Fatalf("expected tenant 123, got %f", tenant)
	}
	if role != "ADMIN" {
		t.Fatalf("expected role 'ADMIN', got %s", role)
	}
	if int64(userID) != 456 {
		t.Fatalf("expected user id 456, got %f", userID)
	}
}

func TestParseJWTInvalid(t *testing.T) {
	_, _, _, _, err := ParseJWT("invalid")
	if err == nil {
		t.Fatalf("expected error parsing invalid token")
	}
}

func TestGenerateAndParseJWTWithBilling(t *testing.T) {
	token, err := GenerateJWTWithBilling("user", 123, "ADMIN", 456, BillingClaims{CanWrite: true})
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	username, tenant, role, userID, billing, err := ParseJWTWithBilling(token)
	if err != nil {
		t.Fatalf("unexpected error parsing token: %v", err)
	}

	if username != "user" {
		t.Fatalf("expected username 'user', got %s", username)
	}
	if int64(tenant) != 123 {
		t.Fatalf("expected tenant 123, got %f", tenant)
	}
	if role != "ADMIN" {
		t.Fatalf("expected role 'ADMIN', got %s", role)
	}
	if int64(userID) != 456 {
		t.Fatalf("expected user id 456, got %f", userID)
	}
	if !billing.CanWrite {
		t.Fatalf("expected billing CanWrite true")
	}
}

func TestParseJWTWithBillingDefaultClaims(t *testing.T) {
	token, err := GenerateJWT("user", 123, "ADMIN", 456)
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	_, _, _, _, billing, err := ParseJWTWithBilling(token)
	if err != nil {
		t.Fatalf("unexpected error parsing token: %v", err)
	}
	if billing.CanWrite {
		t.Fatalf("expected billing CanWrite false when claim is missing")
	}
}

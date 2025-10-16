package helper

import "testing"

func TestGenerateAndParseJWT(t *testing.T) {
	token, err := GenerateJWT("user", 123)
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	username, tenant, err := ParseJWT(token)
	if err != nil {
		t.Fatalf("unexpected error parsing token: %v", err)
	}

	if username != "user" {
		t.Fatalf("expected username 'user', got %s", username)
	}
	if tenant != 123 {
		t.Fatalf("expected tenant 123, got %f", tenant)
	}
}

func TestParseJWTInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for invalid token")
		}
	}()
	ParseJWT("invalid")
}

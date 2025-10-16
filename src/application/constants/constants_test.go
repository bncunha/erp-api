package constants

import "testing"

func TestConstants(t *testing.T) {
	if TENANT_KEY != "tenant_id" {
		t.Fatalf("unexpected tenant key: %s", TENANT_KEY)
	}
	if USERNAME_KEY != "username" {
		t.Fatalf("unexpected username key: %s", USERNAME_KEY)
	}
}

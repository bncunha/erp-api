package domain

import "testing"

func TestNewUser(t *testing.T) {
	user := NewUser(CreateUserParams{
		Username:    "tester",
		Name:        "Test User",
		PhoneNumber: ptr("123"),
		Role:        "ADMIN",
		Email:       "test@example.com",
	})

	if user.Username != "tester" || user.Name != "Test User" || user.Role != "ADMIN" || user.Email != "test@example.com" {
		t.Fatalf("unexpected user returned: %+v", user)
	}
	if user.PhoneNumber == nil || *user.PhoneNumber != "123" {
		t.Fatalf("expected phone to be set, got %+v", user.PhoneNumber)
	}
}

func TestNewUserEmptyPhoneReturnsNil(t *testing.T) {
	user := NewUser(CreateUserParams{
		Username:    "tester",
		Name:        "Test User",
		Role:        "ADMIN",
		Email:       "test@example.com",
		PhoneNumber: ptr(" "),
	})

	if user.PhoneNumber != nil {
		t.Fatalf("expected phone to be nil, got %+v", user.PhoneNumber)
	}
}

func ptr(s string) *string { return &s }

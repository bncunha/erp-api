package domain

import "testing"

func TestNewUser(t *testing.T) {
	user := NewUser(CreateUserParams{
		Username:    "tester",
		Name:        "Test User",
		PhoneNumber: "123",
		Role:        "ADMIN",
		Email:       "test@example.com",
	})

	if user.Username != "tester" || user.Name != "Test User" || user.Role != "ADMIN" || user.Email != "test@example.com" {
		t.Fatalf("unexpected user returned: %+v", user)
	}
}

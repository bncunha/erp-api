package service

import (
	"context"
	"errors"
	"testing"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/domain"
)

func TestAuthServiceLoginSuccess(t *testing.T) {
	userRepo := &stubUserRepository{getByUsername: domain.User{Username: "user", Password: "password", Name: "User", TenantId: 1}}
	service := &authService{userRepository: userRepo}

	output, err := service.Login(context.Background(), request.LoginRequest{Username: "user", Password: "password"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Name != "User" || output.Token == "" {
		t.Fatalf("expected token and name")
	}
}

func TestAuthServiceLoginInvalidPassword(t *testing.T) {
	userRepo := &stubUserRepository{getByUsername: domain.User{Username: "user", Password: "password"}}
	service := &authService{userRepository: userRepo}

	_, err := service.Login(context.Background(), request.LoginRequest{Username: "user", Password: "wrong"})
	if err == nil {
		t.Fatalf("expected error for wrong password")
	}
}

func TestAuthServiceLoginValidationError(t *testing.T) {
	service := &authService{}
	if _, err := service.Login(context.Background(), request.LoginRequest{}); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestAuthServiceLoginRepositoryError(t *testing.T) {
	userRepo := &stubUserRepository{getByUsernameErr: errors.New("fail")}
	service := &authService{userRepository: userRepo}
	if _, err := service.Login(context.Background(), request.LoginRequest{Username: "user", Password: "password"}); err == nil || err.Error() != "fail" {
		t.Fatalf("expected repository error")
	}
}

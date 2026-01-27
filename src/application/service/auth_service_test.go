package service

import (
	"context"
	"errors"
	"testing"

	request "github.com/bncunha/erp-api/src/api/requests"
	helper "github.com/bncunha/erp-api/src/application/helpers"
	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/domain"
)

func TestAuthServiceLoginSuccess(t *testing.T) {
	userRepo := &stubUserRepository{getByUsername: domain.User{Username: "user", Password: "password", Name: "User", TenantId: 1}}
	service := &authService{userRepository: userRepo, encrypt: &stubEncrypto{}, billingService: stubBillingService{}}

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
	service := &authService{userRepository: userRepo, encrypt: &stubEncrypto{}, billingService: stubBillingService{}}

	_, err := service.Login(context.Background(), request.LoginRequest{Username: "user", Password: "wrong"})
	if err == nil {
		t.Fatalf("expected error for wrong password")
	}
}

func TestAuthServiceLoginValidationError(t *testing.T) {
	service := &authService{encrypt: &stubEncrypto{}, billingService: stubBillingService{}}
	if _, err := service.Login(context.Background(), request.LoginRequest{}); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestAuthServiceLoginRepositoryError(t *testing.T) {
	userRepo := &stubUserRepository{getByUsernameErr: errors.New("fail")}
	service := &authService{userRepository: userRepo, encrypt: &stubEncrypto{}, billingService: stubBillingService{}}
	if _, err := service.Login(context.Background(), request.LoginRequest{Username: "user", Password: "password"}); err == nil || err.Error() != "fail" {
		t.Fatalf("expected repository error")
	}
}

func TestAuthServiceLoginTokenGenerationError(t *testing.T) {
	userRepo := &stubUserRepository{getByUsername: domain.User{Username: "user", Password: "secret"}}
	service := &authService{
		userRepository: userRepo,
		encrypt:        &stubEncrypto{},
		billingService: stubBillingService{},
		generateToken: func(username string, tenantID int64, role string, userID int64, billing helper.BillingClaims) (string, error) {
			return "", errors.New("token fail")
		},
	}
	if _, err := service.Login(context.Background(), request.LoginRequest{Username: "user", Password: "secret"}); err == nil || err.Error() != "token fail" {
		t.Fatalf("expected token generation error")
	}
}

func TestAuthServiceLoginBcryptPassword(t *testing.T) {
	userRepo := &stubUserRepository{
		getByUsername: domain.User{Username: "user", Password: "$2a$1234567890123456789012", Name: "User", TenantId: 1},
	}
	service := &authService{userRepository: userRepo, encrypt: &stubEncrypto{compareResp: true}, billingService: stubBillingService{}}

	output, err := service.Login(context.Background(), request.LoginRequest{Username: "user", Password: "password"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Token == "" {
		t.Fatalf("expected token to be generated")
	}
}

type stubBillingService struct{}

func (stubBillingService) GetStatus(ctx context.Context) (output.BillingStatusOutput, error) {
	return output.BillingStatusOutput{PlanName: domain.PlanNameTrial, CanWrite: true}, nil
}

func (stubBillingService) GetSummary(ctx context.Context) (output.BillingSummaryOutput, error) {
	return output.BillingSummaryOutput{}, nil
}

func (stubBillingService) GetPayments(ctx context.Context) ([]output.BillingPaymentOutput, error) {
	return nil, nil
}

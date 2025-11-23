package service

import (
	"context"
	"errors"
	"testing"

	"github.com/bncunha/erp-api/src/application/service/input"
	"github.com/bncunha/erp-api/src/domain"
)

func TestUserTokenServiceCreateSuccess(t *testing.T) {
	repo := &stubUserTokenRepository{
		createID: 1,
		getById:  domain.UserToken{Uuid: "generated-uuid"},
	}
	service := NewUserTokenService(repo, &stubEncrypto{})

	token, err := service.Create(context.Background(), input.CreateUserTokenInput{
		User:      domain.User{Id: 1},
		CreatedBy: domain.User{Id: 2},
		Type:      domain.UserTokenTypeInvite,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.createInput.User.Id != 1 || repo.createInput.Type != domain.UserTokenTypeInvite {
		t.Fatalf("expected repository create to receive user token data")
	}
	if token.Uuid != "generated-uuid" || token.Code == "" {
		t.Fatalf("expected token data to be returned, got %+v", token)
	}
}

func TestUserTokenServiceCreateRepositoryError(t *testing.T) {
	repo := &stubUserTokenRepository{createErr: errors.New("fail")}
	service := NewUserTokenService(repo, &stubEncrypto{})
	if _, err := service.Create(context.Background(), input.CreateUserTokenInput{}); err == nil || err.Error() != "fail" {
		t.Fatalf("expected repository error")
	}
}

func TestUserTokenServiceCreateGetByIdError(t *testing.T) {
	repo := &stubUserTokenRepository{getByIdErr: errors.New("fail")}
	service := NewUserTokenService(repo, &stubEncrypto{})
	if _, err := service.Create(context.Background(), input.CreateUserTokenInput{}); err == nil || err.Error() != "fail" {
		t.Fatalf("expected get by id error")
	}
}

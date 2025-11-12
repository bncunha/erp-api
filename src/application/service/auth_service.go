package service

import (
	"context"
	"errors"

	request "github.com/bncunha/erp-api/src/api/requests"
	helper "github.com/bncunha/erp-api/src/application/helpers"
	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/domain"
)

type AuthService interface {
	Login(ctx context.Context, input request.LoginRequest) (output.LoginOutput, error)
}

type authService struct {
	userRepository domain.UserRepository
	generateToken  func(username string, tenantID int64, role string, userID int64) (string, error)
}

func NewAuthService(userRepository domain.UserRepository) AuthService {
	return &authService{
		userRepository: userRepository,
		generateToken:  helper.GenerateJWT,
	}
}

func (s *authService) Login(ctx context.Context, input request.LoginRequest) (out output.LoginOutput, err error) {
	err = input.Validate()
	if err != nil {
		return out, err
	}

	user, err := s.userRepository.GetByUsername(ctx, input.Username)
	if err != nil {
		return out, err
	}

	if user.Password != input.Password {
		return out, errors.New("senha incorreta")
	}

	tokenFn := s.generateToken
	if tokenFn == nil {
		tokenFn = helper.GenerateJWT
	}

	token, err := tokenFn(user.Username, user.TenantId, user.Role, user.Id)
	if err != nil {
		return out, err
	}

	return output.LoginOutput{Token: token, Name: user.Name}, nil

}

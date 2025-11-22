package service

import (
	"context"
	"errors"
	"strings"

	request "github.com/bncunha/erp-api/src/api/requests"
	helper "github.com/bncunha/erp-api/src/application/helpers"
	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/logs"
)

type AuthService interface {
	Login(ctx context.Context, input request.LoginRequest) (output.LoginOutput, error)
}

type authService struct {
	userRepository domain.UserRepository
	encrypt        domain.Encrypto
	generateToken  func(username string, tenantID int64, role string, userID int64) (string, error)
}

func NewAuthService(userRepository domain.UserRepository, encrypt domain.Encrypto) AuthService {
	return &authService{
		userRepository: userRepository,
		generateToken:  helper.GenerateJWT,
		encrypt:        encrypt,
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

	isBcrypt := strings.HasPrefix(user.Password, "$2") && len(user.Password) > 20
	if isBcrypt {
		if match, err := s.encrypt.Compare(user.Password, input.Password); err != nil || !match {
			return out, errors.New("senha incorreta")
		}
	} else {
		if user.Password != input.Password {
			return out, errors.New("senha incorreta")
		} else {
			hashedPassword, err := s.encrypt.Encrypt(user.Password)
			if err != nil {
				logs.Logger.Errorf("Erro ao atualizar a senha do usuário %d para bcrypt: %v", user.Id, err)
			}
			err = s.userRepository.UpdatePassword(ctx, user, hashedPassword)
			if err != nil {
				logs.Logger.Errorf("Erro ao atualizar usuario %d para bcrypt: %v", user.Id, err)
			}
		}
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

package service

import (
	"context"

	"github.com/bncunha/erp-api/src/application/service/input"
	"github.com/bncunha/erp-api/src/domain"
)

type UserTokenService interface {
	Create(ctx context.Context, input input.CreateUserTokenInput) (domain.UserToken, error)
}

type userTokenService struct {
	userTokenRepository domain.UserTokenRepository
	encrypto            domain.Encrypto
}

func NewUserTokenService(userTokenRepository domain.UserTokenRepository, encrypto domain.Encrypto) UserTokenService {
	return &userTokenService{userTokenRepository, encrypto}
}

func (s *userTokenService) Create(ctx context.Context, input input.CreateUserTokenInput) (domain.UserToken, error) {
	token := domain.NewUserToken(domain.CreateUserTokenParams{
		User:      input.User,
		CreatedBy: input.CreatedBy,
		Type:      input.Type,
	}, s.encrypto)
	id, err := s.userTokenRepository.Create(ctx, token)
	if err != nil {
		return domain.UserToken{}, err
	}
	code := token.Code

	userToken, err := s.userTokenRepository.GetById(ctx, id)
	if err != nil {
		return domain.UserToken{}, err
	}
	userToken.Code = code
	return userToken, nil
}

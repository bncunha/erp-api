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
}

func NewUserTokenService(userTokenRepository domain.UserTokenRepository) UserTokenService {
	return &userTokenService{userTokenRepository}
}

func (s *userTokenService) Create(ctx context.Context, input input.CreateUserTokenInput) (domain.UserToken, error) {
	token := domain.NewUserToken(domain.CreateUserTokenParams{
		User:      input.User,
		CreatedBy: input.CreatedBy,
		Type:      input.Type,
	})
	id, err := s.userTokenRepository.Create(ctx, token)
	if err != nil {
		return domain.UserToken{}, err
	}

	userToken, err := s.userTokenRepository.GetById(ctx, id)
	if err != nil {
		return domain.UserToken{}, err
	}

	return userToken, nil
}

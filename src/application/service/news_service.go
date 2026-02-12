package service

import (
	"context"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
)

type NewsService interface {
	GetLatest(ctx context.Context) (domain.News, error)
}

type newsService struct {
	newsRepository domain.NewsRepository
}

func NewNewsService(newsRepository domain.NewsRepository) NewsService {
	return &newsService{newsRepository: newsRepository}
}

func (s *newsService) GetLatest(ctx context.Context) (domain.News, error) {
	tenantRaw := ctx.Value(constants.TENANT_KEY)
	roleRaw := ctx.Value(constants.ROLE_KEY)

	var tenantId int64
	switch value := tenantRaw.(type) {
	case int64:
		tenantId = value
	case float64:
		tenantId = int64(value)
	default:
		return domain.News{}, errors.New("tenant id invalido")
	}

	role, ok := roleRaw.(string)
	if !ok || role == "" {
		return domain.News{}, errors.New("role invalida")
	}

	return s.newsRepository.GetLatestVisible(ctx, tenantId, domain.Role(role))
}

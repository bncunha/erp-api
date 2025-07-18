package service

import (
	"context"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

type SkuService interface {
	Create(ctx context.Context, request request.CreateSkuRequest, productId int64) (error)
	Update(ctx context.Context, request request.EditSkuRequest, skuId int64) (error)
}

type skuService struct {
	skuRepository repository.SkuRepository
}

func NewSkuService(skuRepository repository.SkuRepository) SkuService {
	return &skuService{skuRepository}
}

func (s *skuService) Create(ctx context.Context, request request.CreateSkuRequest, productId int64) (error) {
	err := request.Validate()
	if err != nil {
		return err
	}

	sku := domain.Sku{
		Code:    request.Code,
		Color:   request.Color,
		Size:    request.Size,
		Cost:    request.Cost,
		Price:   request.Price,
	}

	_, err = s.skuRepository.Create(ctx, sku, productId)
	if err != nil {
		if errors.IsDuplicated(err) {
			return errors.New("Código já existe no banco de dados")
		}
		return err
	}
	return nil
}

func (s *skuService) Update(ctx context.Context, request request.EditSkuRequest, skuId int64) (error) {
	err := request.Validate()
	if err != nil {
		return err
	}

	sku := domain.Sku{
		Id:     skuId,
		Code:    request.Code,
		Color:   request.Color,
		Size:    request.Size,
		Cost:    request.Cost,
		Price:   request.Price,
	}

	err = s.skuRepository.Update(ctx, sku)
	if err != nil {
		if errors.IsDuplicated(err) {
			return errors.New("Código já existe no banco de dados")
		}
		return err
	}
	return nil
}
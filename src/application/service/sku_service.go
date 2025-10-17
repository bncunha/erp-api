package service

import (
	"context"
	"database/sql"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	helper "github.com/bncunha/erp-api/src/application/helpers"
	"github.com/bncunha/erp-api/src/application/service/input"
	"github.com/bncunha/erp-api/src/application/usecase/inventory_usecase"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

type SkuService interface {
	Create(ctx context.Context, request request.CreateSkuRequest, productId int64) error
	Update(ctx context.Context, request request.EditSkuRequest, skuId int64) error
	GetById(ctx context.Context, skuId int64) (domain.Sku, error)
	GetAll(ctx context.Context, filters GetSkusFilters) ([]domain.Sku, error)
	Inactivate(ctx context.Context, id int64) error
}

type skuService struct {
	skuRepository     repository.SkuRepository
	inventoryUseCase  inventory_usecase.InventoryUseCase
	productRepository repository.ProductRepository
	txManager         transactionManager
}

func NewSkuService(skuRepository repository.SkuRepository, inventoryUseCase inventory_usecase.InventoryUseCase, productRepository repository.ProductRepository, txManager transactionManager) SkuService {
	return &skuService{skuRepository, inventoryUseCase, productRepository, txManager}
}

type GetSkusFilters struct {
	SellerId *float64
}

func (s *skuService) Create(ctx context.Context, request request.CreateSkuRequest, productId int64) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	sku := domain.Sku{
		Code:  request.Code,
		Color: request.Color,
		Size:  request.Size,
		Cost:  request.Cost,
		Price: request.Price,
	}

	product, err := s.productRepository.GetById(ctx, productId)
	if err != nil {
		return err
	}

	skuId, err := s.skuRepository.Create(ctx, sku, product.Id)
	if err != nil {
		if errors.IsDuplicated(err) {
			return errors.New("Código já cadastrado!")
		}
		return err
	}

	var tx *sql.Tx
	if s.txManager != nil {
		tx, err = s.txManager.BeginTx(ctx)
		if err != nil {
			return err
		}
		defer func() {
			if err != nil && tx != nil {
				tx.Rollback()
			}
		}()
	}

	if request.Quantity != nil && request.DestinationId != nil {
		err = s.inventoryUseCase.DoTransaction(ctx, tx, inventory_usecase.DoTransactionInput{
			Type:                   domain.InventoryTransactionTypeIn,
			InventoryOriginId:      0,
			InventoryDestinationId: *request.DestinationId,
			Justification:          "Cadastro de Produto",
			Skus:                   []inventory_usecase.DoTransactionSkusInput{{SkuId: skuId, Quantity: *request.Quantity}},
		})
		if err != nil {
			return errors.New("Operação realizada parcialmente! Erro ao atualizar a quantidade de itens no estoque!")
		}
	}
	if tx != nil {
		return tx.Commit()
	}
	return nil
}

func (s *skuService) Update(ctx context.Context, request request.EditSkuRequest, skuId int64) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	sku := domain.Sku{
		Id:    skuId,
		Code:  request.Code,
		Color: request.Color,
		Size:  request.Size,
		Cost:  request.Cost,
		Price: request.Price,
	}

	err = s.skuRepository.Update(ctx, sku)
	if err != nil {
		if errors.IsDuplicated(err) {
			return errors.New("Código já cadastrado!")
		}
		return err
	}
	return nil
}

func (s *skuService) GetById(ctx context.Context, skuId int64) (domain.Sku, error) {
	sku, err := s.skuRepository.GetById(ctx, skuId)
	if err != nil {
		return domain.Sku{}, err
	}
	return sku, nil
}

func (s *skuService) GetAll(ctx context.Context, filters GetSkusFilters) ([]domain.Sku, error) {
	var sellerIdFilter *float64
	if helper.GetRole(ctx) == domain.UserRoleAdmin {
		sellerIdFilter = filters.SellerId
	} else {
		id := ctx.Value(constants.USERID_KEY).(float64)
		sellerIdFilter = &id
	}

	skus, err := s.skuRepository.GetAll(ctx, input.GetSkusInput{SellerId: sellerIdFilter})
	if err != nil {
		return skus, err
	}
	return skus, nil
}

func (s *skuService) Inactivate(ctx context.Context, id int64) error {
	return s.skuRepository.Inactivate(ctx, id)
}

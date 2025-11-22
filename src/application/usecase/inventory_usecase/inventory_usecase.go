package inventory_usecase

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

var (
	ErrEnventoryItemDestinationNotFound = errors.New("Item de estoque não encontrado no destino")
	ErrInventoryItemOriginNotFound      = errors.New("Item de estoque não encontrado na origem")
	ErrQuantityInsufficient             = errors.New("Quantidade insuficiente")
	ErrInventoryesTransferEquais        = errors.New("Inventários de origem e de destino precisam ser diferentes")
	ErrInventoryItemNotFound            = errors.New("Item de estoque não encontrado")
	ErrSkusNotFound                     = errors.New("SKUs não encontrados")
	ErrSkusDuplicated                   = errors.New("SKUs duplicados encontrados")
)

type InventoryUseCase interface {
	DoTransaction(ctx context.Context, tx *sql.Tx, input DoTransactionInput) error
}

type inventoryUseCase struct {
	repository               *repository.Repository
	inventoryRepository      domain.InventoryRepository
	inventoryItemRepository  domain.InventoryItemRepository
	inventoryTransactionRepo domain.InventoryTransactionRepository
	skuRepository            domain.SkuRepository
}

func NewInventoryUseCase(repository *repository.Repository, inventoryRepository domain.InventoryRepository, inventoryItemRepository domain.InventoryItemRepository, inventoryTransactionRepo domain.InventoryTransactionRepository, skuRepository domain.SkuRepository) InventoryUseCase {
	return &inventoryUseCase{repository, inventoryRepository, inventoryItemRepository, inventoryTransactionRepo, skuRepository}
}

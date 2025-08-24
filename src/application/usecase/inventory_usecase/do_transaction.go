package inventory_usecase

import (
	"context"
	"database/sql"
	"time"

	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

var (
	ErrEnventoryItemDestinationNotFound = errors.New("Item de estoque não encontrado no destino")
	ErrInventoryItemOriginNotFound      = errors.New("Item de estoque não encontrado na origem")
	ErrQuantityInsufficient             = errors.New("Quantidade insuficiente")
	ErrInventoryesTransferEquais        = errors.New("Inventários de origem e de destino precisam ser diferentes")
)

type InventoryUseCase interface {
	DoTransaction(ctx context.Context, input DoTransactionInput) error
}

type inventoryUseCase struct {
	repository               *repository.Repository
	inventoryRepository      repository.InventoryRepository
	inventoryItemRepository  repository.InventoryItemRepository
	inventoryTransactionRepo repository.InventoryTransactionRepository
	skuRepository            repository.SkuRepository
}

func NewInventoryUseCase(repository *repository.Repository, inventoryRepository repository.InventoryRepository, inventoryItemRepository repository.InventoryItemRepository, inventoryTransactionRepo repository.InventoryTransactionRepository, skuRepository repository.SkuRepository) InventoryUseCase {
	return &inventoryUseCase{repository, inventoryRepository, inventoryItemRepository, inventoryTransactionRepo, skuRepository}
}

func (s *inventoryUseCase) DoTransaction(ctx context.Context, input DoTransactionInput) (err error) {

	var inventoryOut domain.Inventory
	var inventoryIn domain.Inventory
	var inventoryItemOut domain.InventoryItem
	var inventoryItemIn domain.InventoryItem

	if input.InventoryDestinationId != 0 {
		inventoryIn, err = s.inventoryRepository.GetById(ctx, input.InventoryDestinationId)
		if err != nil {
			return err
		}
	}
	if input.InventoryOriginId != 0 {
		inventoryOut, err = s.inventoryRepository.GetById(ctx, input.InventoryOriginId)
		if err != nil {
			return err
		}
	}

	sku, err := s.skuRepository.GetById(ctx, input.SkuId)
	if err != nil {
		return err
	}

	inventoryItemOut, err = s.inventoryItemRepository.GetBySkuIdAndInventoryId(ctx, sku.Id, inventoryOut.Id)
	if err != nil && !errors.Is(err, repository.ErrInventoryItemNotFound) {
		return err
	}

	inventoryItemIn, err = s.inventoryItemRepository.GetBySkuIdAndInventoryId(ctx, sku.Id, inventoryIn.Id)
	if err != nil && !errors.Is(err, repository.ErrInventoryItemNotFound) {
		return err
	}

	err = s.validateInventoryTransaction(ctx, inventoryItemOut, inventoryItemIn, input.Quantity, input.Type)
	if err != nil {
		return err
	}

	tx, err := s.repository.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	inventoryItemIn, err = s.createInventoryItemInIfNotExists(ctx, tx, sku, input.Quantity, inventoryItemIn, domain.InventoryTransactionType(input.Type), inventoryIn)
	if err != nil {
		return err
	}

	err = s.createTransactions(ctx, tx, inventoryItemOut, inventoryItemIn, domain.InventoryTransaction{
		Quantity:      input.Quantity,
		Date:          time.Now(),
		InventoryOut:  inventoryOut,
		InventoryIn:   inventoryIn,
		Justification: input.Justification,
	}, input.Type)
	if err != nil {
		return err
	}

	switch input.Type {
	case domain.InventoryTransactionTypeTransfer:
		err = s.transferQuantity(ctx, tx, inventoryItemOut, inventoryItemIn, input.Quantity)
		if err != nil {
			return err
		}
	case domain.InventoryTransactionTypeIn:
		err = s.addQuantity(ctx, tx, inventoryItemIn, input.Quantity)
		if err != nil {
			return err
		}
	case domain.InventoryTransactionTypeOut:
		err = s.subQuantity(ctx, tx, inventoryItemOut, input.Quantity)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *inventoryUseCase) createTransactions(ctx context.Context, tx *sql.Tx, inventoryItemOut domain.InventoryItem, inventoryItemIn domain.InventoryItem, transaction domain.InventoryTransaction, transactionType domain.InventoryTransactionType) error {
	if transactionType == domain.InventoryTransactionTypeTransfer {
		transaction.InventoryItem = inventoryItemOut
		transaction.Type = domain.InventoryTransactionTypeTransfer
		_, err := s.inventoryTransactionRepo.Create(ctx, tx, transaction)
		if err != nil {
			return err
		}
	} else if transactionType == domain.InventoryTransactionTypeIn {
		transaction.InventoryItem = inventoryItemIn
		transaction.Type = domain.InventoryTransactionTypeIn
		_, err := s.inventoryTransactionRepo.Create(ctx, tx, transaction)
		if err != nil {
			return err
		}
	} else if transactionType == domain.InventoryTransactionTypeOut {
		transaction.InventoryItem = inventoryItemOut
		transaction.Type = domain.InventoryTransactionTypeOut
		_, err := s.inventoryTransactionRepo.Create(ctx, tx, transaction)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *inventoryUseCase) createInventoryItemInIfNotExists(ctx context.Context, tx *sql.Tx, sku domain.Sku, quantity float64, inventoryItemIn domain.InventoryItem, transactionType domain.InventoryTransactionType, inventoryIn domain.Inventory) (domain.InventoryItem, error) {
	if inventoryItemIn.Id != 0 {
		return inventoryItemIn, nil
	}
	if transactionType == domain.InventoryTransactionTypeIn && inventoryItemIn.Id == 0 {
		insertInventoryItem := domain.InventoryItem{
			InventoryId: inventoryIn.Id,
			SkuId:       sku.Id,
			Quantity:    0,
		}
		inventoryItemId, err := s.inventoryItemRepository.Create(ctx, tx, insertInventoryItem)
		if err != nil {
			return domain.InventoryItem{}, err
		}
		inventoryItemIn, err = s.inventoryItemRepository.GetByIdWithTransaction(ctx, tx, inventoryItemId)
		if err != nil {
			return domain.InventoryItem{}, err
		}
	}
	return inventoryItemIn, nil
}

func (s *inventoryUseCase) validateInventoryTransaction(ctx context.Context, inventoryItemOut domain.InventoryItem, inventoryItemIn domain.InventoryItem, quantity float64, inventoryType domain.InventoryTransactionType) error {
	if inventoryType == domain.InventoryTransactionTypeTransfer || inventoryType == domain.InventoryTransactionTypeOut {
		if inventoryItemOut.Id == 0 {
			return ErrInventoryItemOriginNotFound
		}

		if inventoryItemOut.Quantity < quantity {
			return ErrQuantityInsufficient
		}
	}
	if inventoryType == domain.InventoryTransactionTypeTransfer {
		if inventoryItemIn.Id == inventoryItemOut.Id {
			return ErrInventoryesTransferEquais
		}
	}
	return nil
}

func (s *inventoryUseCase) transferQuantity(ctx context.Context, tx *sql.Tx, inventoryItemOut domain.InventoryItem, inventoryItemIn domain.InventoryItem, quantity float64) error {
	err := s.subQuantity(ctx, tx, inventoryItemOut, quantity)
	if err != nil {
		return err
	}
	err = s.addQuantity(ctx, tx, inventoryItemIn, quantity)
	if err != nil {
		return err
	}
	return nil
}

func (s *inventoryUseCase) addQuantity(ctx context.Context, tx *sql.Tx, inventoryItem domain.InventoryItem, quantity float64) error {
	inventoryItem.Quantity = inventoryItem.Quantity + quantity
	err := s.inventoryItemRepository.UpdateQuantity(ctx, tx, inventoryItem)
	if err != nil {
		return err
	}
	return nil
}

func (s *inventoryUseCase) subQuantity(ctx context.Context, tx *sql.Tx, inventoryItem domain.InventoryItem, quantity float64) error {
	inventoryItem.Quantity = inventoryItem.Quantity - quantity
	if inventoryItem.Quantity < 0 {
		return ErrQuantityInsufficient
	}
	err := s.inventoryItemRepository.UpdateQuantity(ctx, tx, inventoryItem)
	if err != nil {
		return err
	}
	return nil
}

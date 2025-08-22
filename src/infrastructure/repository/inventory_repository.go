package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
)

type InventoryRepository interface {
	Create(ctx context.Context, inventory domain.Inventory) (int64, error)
	GetById(ctx context.Context, id int64) (domain.Inventory, error)
	GetAll(ctx context.Context) ([]domain.Inventory, error)
}

type inventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) InventoryRepository {
	return &inventoryRepository{db}
}

func (r *inventoryRepository) Create(ctx context.Context, inventory domain.Inventory) (int64, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var insertedID int64

	query := `INSERT INTO inventories (user_id, tenant_id, type) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, inventory.User.Id, tenantId, inventory.Type).Scan(&insertedID)
	return insertedID, err
}

func (r *inventoryRepository) GetById(ctx context.Context, id int64) (domain.Inventory, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var inventory domain.Inventory

	query := `SELECT id, tenant_id, type FROM inventories WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, id, tenantId).Scan(&inventory.Id, &inventory.TenantId, &inventory.Type)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return inventory, errors.New("Inventário não encontrado")
		}
		return inventory, err
	}

	return inventory, nil
}

func (r *inventoryRepository) GetAll(ctx context.Context) ([]domain.Inventory, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var inventories []domain.Inventory

	query := `SELECT i.id, u.id, u.name, i.tenant_id, i.type
	FROM inventories i
	LEFT JOIN users u ON u.id = i.user_id
	WHERE i.tenant_id = $1 AND i.deleted_at IS NULL ORDER BY i.id ASC`
	rows, err := r.db.QueryContext(ctx, query, tenantId)
	if err != nil {
		return inventories, err
	}
	defer rows.Close()

	for rows.Next() {
		var inventory domain.Inventory
		var nullableUserId sql.NullInt64
		var userName sql.NullString

		err = rows.Scan(&inventory.Id, &nullableUserId, &userName, &inventory.TenantId, &inventory.Type)
		if err != nil {
			return inventories, err
		}
		if nullableUserId.Valid {
			inventory.User.Id = nullableUserId.Int64
			inventory.User.Name = userName.String
		}
		inventories = append(inventories, inventory)
	}
	return inventories, err
}

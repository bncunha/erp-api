package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
)

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) domain.CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) Create(ctx context.Context, category domain.Category) (int64, error) {
	var insertedID int64

	query := `INSERT INTO categories (name, tenant_id) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, category.Name, ctx.Value(constants.TENANT_KEY)).Scan(&insertedID)
	if err != nil {
		if errors.IsUniqueViolation(err) {
			return insertedID, errors.New("Categoria já cadastrada!")
		}
		return insertedID, err
	}
	return insertedID, nil
}

func (r *categoryRepository) GetById(ctx context.Context, id int64) (domain.Category, error) {
	var category domain.Category

	query := `SELECT id, name FROM categories WHERE id = $1 AND tenant_id = $2`
	err := r.db.QueryRowContext(ctx, query, id, ctx.Value(constants.TENANT_KEY)).Scan(&category.Id, &category.Name)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return category, errors.New("Categoria não encontrada")
		}
		return category, err
	}
	return category, nil
}

func (r *categoryRepository) GetByName(ctx context.Context, name string) (domain.Category, error) {
	var category domain.Category

	query := `SELECT id, name FROM categories WHERE name = $1 AND tenant_id = $2`
	err := r.db.QueryRowContext(ctx, query, name, ctx.Value(constants.TENANT_KEY)).Scan(&category.Id, &category.Name)
	if err != nil {
		if errors.IsNoRowsFinded(err) {
			return category, errors.New("Categoria não encontrada")
		}
		return category, err
	}
	return category, nil
}

func (r *categoryRepository) Update(ctx context.Context, category domain.Category) error {
	tenantId := ctx.Value(constants.TENANT_KEY)
	query := `UPDATE categories SET name = $1 WHERE id = $2 AND tenant_id = $3`
	_, err := r.db.ExecContext(ctx, query, category.Name, category.Id, tenantId)
	return err
}

func (r *categoryRepository) Delete(ctx context.Context, id int64) error {
	tenantId := ctx.Value(constants.TENANT_KEY)
	query := `DELETE FROM categories WHERE id = $1 AND tenant_id = $2`
	result, err := r.db.ExecContext(ctx, query, id, tenantId)
	if err != nil {
		if errors.IsForeignKeyViolation(err) {
			return errors.New("Não é possível deletar a categoria pois existem registros associados.")
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Categoria não encontrada")
	}

	return nil
}

func (r *categoryRepository) GetAll(ctx context.Context) ([]domain.Category, error) {
	tenantId := ctx.Value(constants.TENANT_KEY)
	var categories []domain.Category

	query := `SELECT id, name FROM categories WHERE tenant_id = $1 AND deleted_at IS NULL ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query, tenantId)
	if err != nil {
		return categories, err
	}
	defer rows.Close()

	for rows.Next() {
		var category domain.Category
		err = rows.Scan(&category.Id, &category.Name)
		if err != nil {
			return categories, err
		}
		categories = append(categories, category)
	}
	return categories, err
}

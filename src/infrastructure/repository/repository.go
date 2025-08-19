package repository

import "database/sql"

type Repository struct {
	db                  *sql.DB
	ProductRepository   ProductRepository
	CategoryRepository  CategoryRepository
	SkuRepository       SkuRepository
	UserRepository      UserRepository
	InventoryRepository InventoryRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SetupRepositories() {
	r.ProductRepository = NewProductRepository(r.db)
	r.CategoryRepository = NewCategoryRepository(r.db)
	r.SkuRepository = NewSkuRepository(r.db)
	r.UserRepository = NewUserRepository(r.db)
	r.InventoryRepository = NewInventoryRepository(r.db)
}

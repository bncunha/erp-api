package repository

import "database/sql"

type Repository struct {
	db                *sql.DB
	ProductRepository ProductRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SetupRepositories() {
	r.ProductRepository = NewProductRepository(r.db)
}
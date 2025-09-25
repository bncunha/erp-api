package service

import (
	"context"

	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

type CustomerService interface {
	GetAll(ctx context.Context) ([]domain.Customer, error)
}

type customerService struct {
	customerRepository repository.CustomerRepository
}

func NewCustomerService(customerRepository repository.CustomerRepository) CustomerService {
	return &customerService{customerRepository}
}

func (s *customerService) GetAll(ctx context.Context) ([]domain.Customer, error) {
	customers, err := s.customerRepository.GetAll(ctx)
	if err != nil {
		return customers, err
	}
	return customers, nil
}

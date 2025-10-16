package service

import (
	"context"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

type CustomerService interface {
	Create(ctx context.Context, input request.CreateCustomerRequest) (int64, error)
	GetAll(ctx context.Context) ([]domain.Customer, error)
}

type customerService struct {
	customerRepository repository.CustomerRepository
}

func NewCustomerService(customerRepository repository.CustomerRepository) CustomerService {
	return &customerService{customerRepository}
}

func (s *customerService) Create(ctx context.Context, input request.CreateCustomerRequest) (int64, error) {
	if err := input.Validate(); err != nil {
		return 0, err
	}
	return s.customerRepository.Create(ctx, domain.Customer{
		Name:        input.Name,
		PhoneNumber: input.Cellphone,
	})
}

func (s *customerService) GetAll(ctx context.Context) ([]domain.Customer, error) {
	customers, err := s.customerRepository.GetAll(ctx)
	if err != nil {
		return customers, err
	}
	return customers, nil
}

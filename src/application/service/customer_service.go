package service

import (
	"context"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/domain"
)

type CustomerService interface {
	Create(ctx context.Context, input request.CreateCustomerRequest) (int64, error)
	GetAll(ctx context.Context) ([]domain.Customer, error)
	GetById(ctx context.Context, id int64) (domain.Customer, error)
	Edit(ctx context.Context, input request.EditCustomerRequest) error
	Inactivate(ctx context.Context, id int64) error
}

type customerService struct {
	customerRepository domain.CustomerRepository
}

func NewCustomerService(customerRepository domain.CustomerRepository) CustomerService {
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

func (s *customerService) GetById(ctx context.Context, id int64) (domain.Customer, error) {
	customer, err := s.customerRepository.GetById(ctx, id)
	if err != nil {
		return customer, err
	}
	return customer, nil
}

func (s *customerService) Edit(ctx context.Context, input request.EditCustomerRequest) error {
	if err := input.Validate(); err != nil {
		return err
	}
	_, err := s.customerRepository.Edit(ctx, domain.Customer{
		Name:        input.Name,
		PhoneNumber: input.Cellphone,
	}, input.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s *customerService) Inactivate(ctx context.Context, id int64) error {
	return s.customerRepository.Inactivate(ctx, id)
}

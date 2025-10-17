package service

import (
	"context"
	"errors"
	"testing"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/domain"
)

func TestCustomerServiceCreate(t *testing.T) {
	repo := &stubCustomerRepository{}
	service := NewCustomerService(repo)

	id, err := service.Create(context.Background(), request.CreateCustomerRequest{Name: "Alice", Cellphone: "123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == 0 {
		t.Fatalf("expected id to be set")
	}
	if repo.created.Name != "Alice" || repo.created.PhoneNumber != "123" {
		t.Fatalf("unexpected customer saved: %+v", repo.created)
	}
}

func TestCustomerServiceCreateValidationError(t *testing.T) {
	service := NewCustomerService(&stubCustomerRepository{})
	if _, err := service.Create(context.Background(), request.CreateCustomerRequest{}); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestCustomerServiceCreateRepositoryError(t *testing.T) {
	expected := errors.New("fail")
	repo := &stubCustomerRepository{createErr: expected}
	service := NewCustomerService(repo)
	if _, err := service.Create(context.Background(), request.CreateCustomerRequest{Name: "Bob", Cellphone: "321"}); err != expected {
		t.Fatalf("expected %v, got %v", expected, err)
	}
}

func TestCustomerServiceGetAll(t *testing.T) {
	repo := &stubCustomerRepository{getAll: []domain.Customer{{Id: 1}}}
	service := NewCustomerService(repo)
	customers, err := service.GetAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(customers) != 1 || customers[0].Id != 1 {
		t.Fatalf("unexpected customers: %+v", customers)
	}

	repo = &stubCustomerRepository{getAllErr: errors.New("fail")}
	service = NewCustomerService(repo)
	if _, err := service.GetAll(context.Background()); err == nil {
		t.Fatalf("expected error")
	}
}

func TestCustomerServiceGetById(t *testing.T) {
	repo := &stubCustomerRepository{getById: domain.Customer{Id: 5}}
	service := NewCustomerService(repo)
	customer, err := service.GetById(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if customer.Id != 5 {
		t.Fatalf("unexpected customer: %+v", customer)
	}

	repo = &stubCustomerRepository{getByIdErr: errors.New("fail")}
	service = NewCustomerService(repo)
	if _, err := service.GetById(context.Background(), 5); err == nil {
		t.Fatalf("expected error")
	}
}

func TestCustomerServiceEdit(t *testing.T) {
	repo := &stubCustomerRepository{}
	service := NewCustomerService(repo)
	err := service.Edit(context.Background(), request.EditCustomerRequest{Id: 1, Name: "Carol", Cellphone: "999"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.created.Name != "Carol" || repo.created.PhoneNumber != "999" {
		t.Fatalf("unexpected updated customer: %+v", repo.created)
	}

	if err := service.Edit(context.Background(), request.EditCustomerRequest{}); err == nil {
		t.Fatalf("expected validation error")
	}

	repo = &stubCustomerRepository{editErr: errors.New("fail")}
	service = NewCustomerService(repo)
	if err := service.Edit(context.Background(), request.EditCustomerRequest{Id: 1, Name: "Eve", Cellphone: "000"}); err == nil {
		t.Fatalf("expected repository error")
	}
}

func TestCustomerServiceInactivate(t *testing.T) {
	repo := &stubCustomerRepository{}
	service := NewCustomerService(repo)
	if err := service.Inactivate(context.Background(), 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	repo = &stubCustomerRepository{inactivateErr: errors.New("fail")}
	service = NewCustomerService(repo)
	if err := service.Inactivate(context.Background(), 3); err == nil {
		t.Fatalf("expected error")
	}
}

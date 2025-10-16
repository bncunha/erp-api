package service

import (
	"context"
	"errors"
	"testing"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

func TestUserServiceCreateReseller(t *testing.T) {
	userRepo := &stubUserRepository{}
	inventoryRepo := &stubInventoryRepository{}
	service := &userService{userRepository: userRepo, inventoryRepository: inventoryRepo}

	req := request.CreateUserRequest{Username: "user", Name: "User", Password: "password", Role: string(domain.InventoryTypeReseller)}
	if err := service.Create(context.Background(), req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inventoryRepo.created.Type != domain.InventoryTypeReseller {
		t.Fatalf("expected reseller inventory to be created")
	}
}

func TestUserServiceCreateDuplicated(t *testing.T) {
	userRepo := &stubUserRepository{createErr: errors.New("duplicate key value violates unique constraint")}
	service := &userService{userRepository: userRepo, inventoryRepository: &stubInventoryRepository{}}

	req := request.CreateUserRequest{Username: "user", Name: "User", Password: "password", Role: "ADMIN"}
	err := service.Create(context.Background(), req)
	if err == nil || err.Error() != "Usuário já cadastrado!" {
		t.Fatalf("expected duplicated error")
	}
}

func TestUserServiceCreateRepositoryError(t *testing.T) {
	userRepo := &stubUserRepository{createErr: errors.New("fail")}
	service := &userService{userRepository: userRepo, inventoryRepository: &stubInventoryRepository{}}
	req := request.CreateUserRequest{Username: "user", Name: "User", Password: "password", Role: "ADMIN"}
	if err := service.Create(context.Background(), req); err == nil || err.Error() != "fail" {
		t.Fatalf("expected repository error")
	}
}

func TestUserServiceUpdateResellerCreatesInventory(t *testing.T) {
	userRepo := &stubUserRepository{}
	inventoryRepo := &stubInventoryRepository{getByUserErr: repository.ErrInventoryNotFound}
	service := &userService{userRepository: userRepo, inventoryRepository: inventoryRepo}

	req := request.EditUserRequest{Username: "user", Name: "User", Role: string(domain.InventoryTypeReseller)}
	if err := service.Update(context.Background(), req, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inventoryRepo.created.Type != domain.InventoryTypeReseller {
		t.Fatalf("expected inventory creation on update")
	}
}

func TestUserServiceUpdateDuplicated(t *testing.T) {
	userRepo := &stubUserRepository{updateErr: errors.New("duplicate key value violates unique constraint")}
	service := &userService{userRepository: userRepo, inventoryRepository: &stubInventoryRepository{}}
	req := request.EditUserRequest{Username: "user", Name: "User", Role: "ADMIN"}

	err := service.Update(context.Background(), req, 1)
	if err == nil || err.Error() != "Usuário já cadastrado!" {
		t.Fatalf("expected duplicated error")
	}
}

func TestUserServiceUpdateExistingInventory(t *testing.T) {
	inventoryRepo := &stubInventoryRepository{getByUser: domain.Inventory{Id: 1}}
	service := &userService{userRepository: &stubUserRepository{}, inventoryRepository: inventoryRepo}
	req := request.EditUserRequest{Username: "user", Name: "User", Role: string(domain.InventoryTypeReseller)}
	if err := service.Update(context.Background(), req, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inventoryRepo.created.Id != 0 {
		t.Fatalf("expected no new inventory creation")
	}
}

func TestUserServiceCreateNonReseller(t *testing.T) {
	inventoryRepo := &stubInventoryRepository{}
	service := &userService{userRepository: &stubUserRepository{}, inventoryRepository: inventoryRepo}
	req := request.CreateUserRequest{Username: "user", Name: "User", Password: "password", Role: "ADMIN"}
	if err := service.Create(context.Background(), req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inventoryRepo.created.Id != 0 {
		t.Fatalf("expected no inventory creation")
	}
}

func TestUserServiceGetters(t *testing.T) {
	userRepo := &stubUserRepository{getById: domain.User{Id: 1}, getAll: []domain.User{{Id: 2}}}
	service := &userService{userRepository: userRepo, inventoryRepository: &stubInventoryRepository{}}

	user, err := service.GetById(context.Background(), 1)
	if err != nil || user.Id != 1 {
		t.Fatalf("unexpected get by id result")
	}

	users, err := service.GetAll(context.Background())
	if err != nil || len(users) != 1 {
		t.Fatalf("unexpected get all result")
	}
}

func TestUserServiceInactivate(t *testing.T) {
	userRepo := &stubUserRepository{}
	service := &userService{userRepository: userRepo, inventoryRepository: &stubInventoryRepository{}}

	if err := service.Inactivate(context.Background(), 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserServiceCreateValidationError(t *testing.T) {
	service := &userService{}
	if err := service.Create(context.Background(), request.CreateUserRequest{}); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestUserServiceCreateInventoryError(t *testing.T) {
	inventoryRepo := &stubInventoryRepository{createErr: errors.New("fail")}
	service := &userService{userRepository: &stubUserRepository{}, inventoryRepository: inventoryRepo}
	req := request.CreateUserRequest{Username: "user", Name: "User", Password: "password", Role: string(domain.InventoryTypeReseller)}

	if err := service.Create(context.Background(), req); err == nil || err.Error() != "fail" {
		t.Fatalf("expected inventory error")
	}
}

func TestUserServiceUpdateValidationError(t *testing.T) {
	service := &userService{}
	if err := service.Update(context.Background(), request.EditUserRequest{}, 1); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestUserServiceUpdateInventoryError(t *testing.T) {
	inventoryRepo := &stubInventoryRepository{getByUserErr: errors.New("fail")}
	service := &userService{userRepository: &stubUserRepository{}, inventoryRepository: inventoryRepo}
	req := request.EditUserRequest{Username: "user", Name: "User", Role: string(domain.InventoryTypeReseller)}

	if err := service.Update(context.Background(), req, 1); err == nil || err.Error() != "fail" {
		t.Fatalf("expected inventory error")
	}
}

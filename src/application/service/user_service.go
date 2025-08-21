package service

import (
	"context"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

type UserService interface {
	Create(ctx context.Context, request request.CreateUserRequest) error
	Update(ctx context.Context, request request.EditUserRequest, userId int64) error
	GetById(ctx context.Context, userId int64) (domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	Inactivate(ctx context.Context, id int64) error
}

type userService struct {
	userRepository      repository.UserRepository
	inventoryRepository repository.InventoryRepository
}

func NewUserService(userRepository repository.UserRepository, inventoryRepository repository.InventoryRepository) UserService {
	return &userService{userRepository, inventoryRepository}
}

func (s *userService) Create(ctx context.Context, request request.CreateUserRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	user := domain.User{
		Username:    request.Username,
		Name:        request.Name,
		PhoneNumber: request.PhoneNumber,
		Password:    request.Password,
		Role:        request.Role,
	}

	userId, err := s.userRepository.Create(ctx, user)
	if err != nil {
		if errors.IsDuplicated(err) {
			return errors.New("Usuário já cadastrado!")
		}
		return err
	}

	if *request.HasInventory {
		inventory := domain.Inventory{
			UserId: userId,
			Type:   domain.InventoryTypeReseller,
		}
		_, err = s.inventoryRepository.Create(ctx, inventory)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *userService) Update(ctx context.Context, request request.EditUserRequest, userId int64) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	user := domain.User{
		Id:          userId,
		Username:    request.Username,
		Name:        request.Name,
		PhoneNumber: request.PhoneNumber,
		Password:    request.Password,
		Role:        request.Role,
	}

	err = s.userRepository.Update(ctx, user)
	if err != nil {
		if errors.IsDuplicated(err) {
			return errors.New("Usuário já cadastrado!")
		}
		return err
	}
	return nil
}

func (s *userService) GetById(ctx context.Context, userId int64) (domain.User, error) {
	user, err := s.userRepository.GetById(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s *userService) GetAll(ctx context.Context) ([]domain.User, error) {
	users, err := s.userRepository.GetAll(ctx)
	if err != nil {
		return users, err
	}
	return users, nil
}

func (s *userService) Inactivate(ctx context.Context, id int64) error {
	return s.userRepository.Inactivate(ctx, id)
}

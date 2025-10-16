package usecase

import (
	"testing"

	"github.com/bncunha/erp-api/src/infrastructure/repository"
)

func TestNewApplicationUseCase(t *testing.T) {
	repos := &repository.Repository{}
	uc := NewApplicationUseCase(repos)
	if uc.repositories != repos {
		t.Fatalf("expected repositories to be set")
	}
}

func TestSetupUseCases(t *testing.T) {
	repos := &repository.Repository{}
	uc := NewApplicationUseCase(repos)
	uc.SetupUseCases()
	if uc.InventoryUseCase == nil {
		t.Fatalf("expected inventory use case to be initialized")
	}
}

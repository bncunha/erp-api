package usecase

import (
	"testing"

	"github.com/bncunha/erp-api/src/application/ports"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
	config "github.com/bncunha/erp-api/src/main"
)

type fakeEncrypto struct{}

func (fakeEncrypto) Encrypt(text string) (string, error) { return text, nil }
func (fakeEncrypto) Compare(hash string, text string) (bool, error) {
	return hash == text, nil
}

type fakeEmailPort struct{}

func (fakeEmailPort) Send(senderEmail string, senderName string, toEmail string, toName string, subject string, body string) error {
	return nil
}

func newTestPorts() *ports.Ports {
	return ports.NewPorts(fakeEncrypto{}, fakeEmailPort{})
}

func TestNewApplicationUseCase(t *testing.T) {
	repos := &repository.Repository{}
	cfg := &config.Config{}
	uc := NewApplicationUseCase(repos, cfg, newTestPorts())
	if uc.repositories != repos {
		t.Fatalf("expected repositories to be set")
	}
}

func TestSetupUseCases(t *testing.T) {
	repos := &repository.Repository{}
	cfg := &config.Config{}
	uc := NewApplicationUseCase(repos, cfg, newTestPorts())
	uc.SetupUseCases()
	if uc.InventoryUseCase == nil {
		t.Fatalf("expected inventory use case to be initialized")
	}
}

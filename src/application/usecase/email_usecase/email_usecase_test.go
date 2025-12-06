package emailusecase

import (
	"context"
	"strings"
	"testing"

	"github.com/bncunha/erp-api/src/domain"
	config "github.com/bncunha/erp-api/src/main"
)

type stubEmailPort struct {
	senderEmail string
	senderName  string
	toEmail     string
	toName      string
	subject     string
	body        string
	err         error
}

func (s *stubEmailPort) Send(senderEmail string, senderName string, toEmail string, toName string, subject string, body string) error {
	s.senderEmail = senderEmail
	s.senderName = senderName
	s.toEmail = toEmail
	s.toName = toName
	s.subject = subject
	s.body = body
	return s.err
}

func TestEmailUseCaseSendInvite(t *testing.T) {
	cfg := &config.Config{FRONTEND_URL: "http://frontend"}
	emailPort := &stubEmailPort{}
	usecase := NewEmailUseCase(cfg, emailPort)

	user := domain.User{Name: "Tester", Email: "tester@example.com"}
	if err := usecase.SendInvite(context.Background(), user, "code123", "uuid456"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedLink := "http://frontend/redefinir-senha?code=code123&uuid=uuid456"
	if emailPort.toEmail != user.Email || emailPort.toName != user.Name || emailPort.subject != InviteSubject {
		t.Fatalf("expected invite email to be sent to user, got %+v", emailPort)
	}
	if !strings.Contains(emailPort.body, expectedLink) {
		t.Fatalf("expected body to contain link %s, got %s", expectedLink, emailPort.body)
	}
}

func TestEmailUseCaseSendRecoverPassword(t *testing.T) {
	cfg := &config.Config{FRONTEND_URL: "http://frontend"}
	port := &stubEmailPort{}
	usecase := NewEmailUseCase(cfg, port)
	user := domain.User{Name: "User", Email: "user@example.com"}
	if err := usecase.SendRecoverPassword(context.Background(), user, "code", "uuid"); err != nil {
		t.Fatalf("expected recover email to succeed, got %v", err)
	}
	if !strings.Contains(port.body, "code=code&uuid=uuid") {
		t.Fatalf("expected recover link in body, got %s", port.body)
	}
}

func TestEmailUseCaseSendWelcome(t *testing.T) {
	cfg := &config.Config{FRONTEND_URL: "http://frontend"}
	port := &stubEmailPort{}
	usecase := NewEmailUseCase(cfg, port)

	if err := usecase.SendWelcome(context.Background(), "user@test.com", "User Test"); err != nil {
		t.Fatalf("unexpected error sending welcome: %v", err)
	}

	if port.subject != WelcomeSubject {
		t.Fatalf("expected welcome subject, got %s", port.subject)
	}
	if !strings.Contains(port.body, "User Test") {
		t.Fatalf("expected name in body, got %s", port.body)
	}
	if !strings.Contains(port.body, cfg.FRONTEND_URL) {
		t.Fatalf("expected frontend url in body, got %s", port.body)
	}
}

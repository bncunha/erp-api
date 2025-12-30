package errors

import (
	stdErrors "errors"
	"strings"
	"testing"

	"github.com/lib/pq"
)

func TestNew(t *testing.T) {
	err := New("message")
	if err == nil || err.Error() != "message" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIsUniqueViolation(t *testing.T) {
	pqErr := &pq.Error{Code: "23505"}
	if !IsUniqueViolation(pqErr) {
		t.Fatalf("expected true for unique violation")
	}
	if IsUniqueViolation(stdErrors.New("other")) {
		t.Fatalf("expected false for non pq error")
	}
}

func TestIsForeignKeyViolation(t *testing.T) {
	pqErr := &pq.Error{Code: "23503"}
	if !IsForeignKeyViolation(pqErr) {
		t.Fatalf("expected true for foreign key violation")
	}
	if IsForeignKeyViolation(stdErrors.New("other")) {
		t.Fatalf("expected false for non pq error")
	}
}

func TestIsNoRowsFinded(t *testing.T) {
	if !IsNoRowsFinded(stdErrors.New("no rows in result set")) {
		t.Fatalf("expected true")
	}
	if IsNoRowsFinded(stdErrors.New("other")) {
		t.Fatalf("expected false")
	}
}

func TestIsDuplicated(t *testing.T) {
	if !IsDuplicated(stdErrors.New("duplicate key value violates unique constraint")) {
		t.Fatalf("expected true")
	}
	if IsDuplicated(stdErrors.New("other")) {
		t.Fatalf("expected false")
	}
}

func TestIs(t *testing.T) {
	target := stdErrors.New("target")
	err := stdErrors.Join(target)
	if !Is(err, target) {
		t.Fatalf("expected Is to match target")
	}
}

func TestParseDuplicatedMessage(t *testing.T) {
	err := ParseDuplicatedMessage("Usuario", stdErrors.New("other"))
	if !strings.Contains(err.Error(), "Usuario") || !strings.Contains(err.Error(), "cadastrado") {
		t.Fatalf("unexpected error message: %v", err)
	}

	phoneErr := ParseDuplicatedMessage("Usuario", &pq.Error{Detail: "phone_number"})
	if !strings.Contains(phoneErr.Error(), "telefone") {
		t.Fatalf("expected telefone message, got %v", phoneErr)
	}

	emailErr := ParseDuplicatedMessage("Usuario", &pq.Error{Detail: "email"})
	if !strings.Contains(emailErr.Error(), "email") {
		t.Fatalf("expected email message, got %v", emailErr)
	}
}

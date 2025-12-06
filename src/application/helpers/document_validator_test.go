package helper

import "testing"

func TestIsValidCPF(t *testing.T) {
	if !IsValidCPF("390.533.447-05") {
		t.Fatalf("expected valid cpf")
	}
	if IsValidCPF("111.111.111-11") {
		t.Fatalf("expected invalid cpf")
	}
}

func TestIsValidCNPJ(t *testing.T) {
	if !IsValidCNPJ("04.252.011/0001-10") {
		t.Fatalf("expected valid cnpj")
	}
	if IsValidCNPJ("11.111.111/1111-11") {
		t.Fatalf("expected invalid cnpj")
	}
}

func TestSanitizeDocument(t *testing.T) {
	sanitized := SanitizeDocument("04.252.011/0001-10")
	if sanitized != "04252011000110" {
		t.Fatalf("unexpected sanitize result: %s", sanitized)
	}
}

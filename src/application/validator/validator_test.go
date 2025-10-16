package validator

import "testing"

type sample struct {
	Name string `validate:"required"`
}

func TestValidate(t *testing.T) {
	t.Run("valid struct", func(t *testing.T) {
		s := sample{Name: "valid"}
		if err := Validate(s); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("invalid struct", func(t *testing.T) {
		s := sample{}
		if err := Validate(s); err == nil {
			t.Fatalf("expected validation error")
		}
	})
}

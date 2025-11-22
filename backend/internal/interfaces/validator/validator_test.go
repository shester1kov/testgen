package validator

import (
	"testing"
)

type sample struct {
	Email string `validate:"required,email"`
	Name  string `validate:"required,min=3"`
}

func TestValidator_Validate(t *testing.T) {
	v := NewValidator()

	if err := v.Validate(sample{}); err == nil {
		t.Fatalf("expected validation errors")
	}

	if err := v.Validate(sample{Email: "user@example.com", Name: "Bob"}); err != nil {
		t.Fatalf("expected valid struct, got %v", err)
	}
}

func TestValidator_ValidateVar(t *testing.T) {
	v := NewValidator()

	if err := v.ValidateVar("not-an-email", "email"); err == nil {
		t.Fatalf("expected email validation failure")
	}

	if err := v.ValidateVar("user@example.com", "email"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := v.ValidateVar("maybe", "oneof=yes no"); err == nil {
		t.Fatalf("expected oneof validation failure")
	}

	if err := v.ValidateVar("not-a-uuid", "uuid"); err == nil {
		t.Fatalf("expected uuid validation failure")
	}
}

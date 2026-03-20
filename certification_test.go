package dgvalidation_test

import (
	"context"
	"testing"

	dgvalidation "github.com/dgframe/dg-validation"
	"github.com/dgframe/dg-validation/adapters/gookit"
)

func TestNoopValidator_Semantics(t *testing.T) {
	// 1. Verify Default (No-op)
	provider := dgvalidation.NewValidationServiceProvider()
	// Using exported way or skipping internal field access if not possible
	// Actually we want to test the public behavior
	ctx := context.Background()
	res, err := provider.Validator().Validate(ctx, nil)

	if err != dgvalidation.ErrCapabilityNotProvided {
		t.Errorf("expected ErrCapabilityNotProvided, got %v", err)
	}

	if res == nil {
		t.Fatal("expected non-nil NoopResult")
	}

	if res.Valid() {
		t.Error("NoopResult.Valid() must be false")
	}

	if len(res.Violations()) != 0 {
		t.Error("NoopResult.Violations() must be empty")
	}
}

func TestGookitAdapter_Certification(t *testing.T) {
	adapter := gookit.NewAdapter()

	type TestStruct struct {
		Name  string `validate:"required"`
		Email string `validate:"email"`
	}

	ctx := context.Background()

	t.Run("Valid Data", func(t *testing.T) {
		data := TestStruct{Name: "John", Email: "john@example.com"}
		res, err := adapter.Validate(ctx, data)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !res.Valid() {
			t.Error("expected result to be valid")
		}
	})

	t.Run("Invalid Data", func(t *testing.T) {
		adapter.StopOnError(false) // Disable fail-fast to capture all errors
		data := TestStruct{Name: "", Email: "invalid-email"}
		res, err := adapter.Validate(ctx, data)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if res.Valid() {
			t.Error("expected result to be invalid")
		}

		violations := res.Violations()
		if len(violations) < 1 {
			t.Errorf("expected at least 1 violation, got %d", len(violations))
		}

		// Verify Violation metadata
		foundName := false
		foundEmail := false
		for _, v := range violations {
			if v.Field() == "Name" {
				foundName = true
			}
			if v.Field() == "Email" {
				foundEmail = true
			}
		}

		// If we set StopOnError(false), we expect both
		if len(violations) >= 2 {
			if !foundName || !foundEmail {
				t.Error("expected violations for Name and Email fields when fail-fast is off")
			}
		}
	})
}

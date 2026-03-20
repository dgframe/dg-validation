package dgvalidation

import (
	"context"
)

// Subject represents the opaque data being validated.
type Subject interface{}

// Validator defines the authoritative contract for validation execution.
type Validator interface {
	// Validate checks a Subject against its defined rules.
	// It returns a Result, where an empty result indicates success.
	Validate(ctx context.Context, data Subject, scene ...string) (Result, error)
}

// Result represents the outcome of a validation pass.
type Result interface {
	// Valid returns true if there are no violations.
	Valid() bool

	// Violations returns all rule failures found during validation.
	Violations() []Violation

	// ByField returns all violations grouped by field path.
	ByField() map[string][]Violation

	// First returns the first violation found, or nil.
	First() Violation
}

// Violation represents a single validation rule failure.
type Violation interface {
	// Field returns the logical path to the invalid field (e.g., "user.email").
	Field() string

	// Rule returns the identity of the failed rule (e.g., "required", "email").
	Rule() string

	// Message returns a generic, non-localized description of the failure.
	Message() string

	// Value returns the actual value that failed validation.
	Value() any
}

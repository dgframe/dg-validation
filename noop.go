package dgvalidation

import (
	"context"
	"errors"
)

// ErrCapabilityNotProvided is returned when the validator capability is missing.
var ErrCapabilityNotProvided = errors.New("capability not provided: validator")

// NoopValidator is a fail-fast Validator used when no adapter is injected.
// It enforces Profile A2 (Hard Cross-Cutting) failure semantics.
type NoopValidator struct{}

// NewNoopValidator creates a new fail-fast NoopValidator.
func NewNoopValidator() *NoopValidator {
	return &NoopValidator{}
}

func (n *NoopValidator) Validate(ctx context.Context, data Subject, scene ...string) (Result, error) {
	return &NoopResult{}, ErrCapabilityNotProvided
}

// NoopResult is a fail-fast Result returned by NoopValidator.
type NoopResult struct{}

func (r *NoopResult) Valid() bool                     { return false }
func (r *NoopResult) Violations() []Violation         { return nil }
func (r *NoopResult) ByField() map[string][]Violation { return nil }
func (r *NoopResult) First() Violation                { return nil }

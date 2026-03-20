package dgvalidation

import (
	"github.com/dgframe/core/foundation"
)

// ValidatorBinding is the typed binding key for the Validator capability.
const ValidatorBinding = Binding

// ValidationServiceProvider is the capability surface provider for validation.
type ValidationServiceProvider struct {
	validator Validator
}

// NewValidationServiceProvider creates a new validation provider.
func NewValidationServiceProvider() *ValidationServiceProvider {
	return &ValidationServiceProvider{
		validator: NewNoopValidator(), // Defaults to Fail-Fast
	}
}

// WithValidator allows injecting a concrete validator implementation.
func (p *ValidationServiceProvider) WithValidator(v Validator) *ValidationServiceProvider {
	p.validator = v
	return p
}

// Validator returns the currently configured validator.
func (p *ValidationServiceProvider) Validator() Validator {
	return p.validator
}

// Register registers the validator capability into the container.
func (p *ValidationServiceProvider) Register(r foundation.Registrar) error {
	r.Bind(ValidatorBinding, func() (interface{}, error) {
		return p.validator, nil
	})
	return nil
}

// Boot boots the framework services.
func (p *ValidationServiceProvider) Boot(ctx foundation.RuntimeContext) error {
	// No-op
	return nil
}

// Shutdown performs no-op shutdown logic.
func (p *ValidationServiceProvider) Shutdown(ctx foundation.RuntimeContext) error {
	// No-op
	return nil
}

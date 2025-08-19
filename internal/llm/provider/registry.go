package provider

import (
	"fmt"
	"sync"

	"github.com/charmbracelet/catwalk/pkg/catwalk"
	"github.com/charmbracelet/crush/internal/config"
)

// ProviderConstructor is a function that creates a provider client from options
type ProviderConstructor func(opts providerClientOptions) ProviderClient

// ProviderRegistration holds information about a registered provider
type ProviderRegistration struct {
	ID          string
	Name        string
	Constructor ProviderConstructor
	Type        catwalk.Type
	SupportsOAuth bool
}

// registryType implements the OAuth registry interface
type registryType struct {
	sync.RWMutex
	providers map[string]*ProviderRegistration
}

// providerRegistry holds all dynamically registered providers
var providerRegistry = &registryType{
	providers: make(map[string]*ProviderRegistration),
}

// RegisterProvider registers a custom provider that can be used by the factory
func RegisterProvider(registration *ProviderRegistration) error {
	if registration == nil {
		return fmt.Errorf("provider registration cannot be nil")
	}
	
	if registration.ID == "" {
		return fmt.Errorf("provider ID cannot be empty")
	}
	
	if registration.Constructor == nil {
		return fmt.Errorf("provider constructor cannot be nil")
	}
	
	providerRegistry.Lock()
	defer providerRegistry.Unlock()
	
	if _, exists := providerRegistry.providers[registration.ID]; exists {
		return fmt.Errorf("provider with ID %s is already registered", registration.ID)
	}
	
	providerRegistry.providers[registration.ID] = registration
	return nil
}

// GetRegisteredProvider retrieves a registered provider by ID
func GetRegisteredProvider(id string) (*ProviderRegistration, bool) {
	providerRegistry.RLock()
	defer providerRegistry.RUnlock()
	
	registration, exists := providerRegistry.providers[id]
	return registration, exists
}

// ListRegisteredProviders returns all registered provider IDs
func ListRegisteredProviders() []string {
	providerRegistry.RLock()
	defer providerRegistry.RUnlock()
	
	var ids []string
	for id := range providerRegistry.providers {
		ids = append(ids, id)
	}
	return ids
}

// IsRegisteredProvider checks if a provider ID is registered
func IsRegisteredProvider(id string) bool {
	_, exists := GetRegisteredProvider(id)
	return exists
}

// IsRegisteredOAuthProvider checks if a registered provider supports OAuth
func IsRegisteredOAuthProvider(id string) bool {
	registration, exists := GetRegisteredProvider(id)
	if !exists {
		return false
	}
	return registration.SupportsOAuth
}

// CreateFromRegistry creates a provider using the registry
func CreateFromRegistry(cfg config.ProviderConfig, opts providerClientOptions) (ProviderClient, bool) {
	registration, exists := GetRegisteredProvider(cfg.ID)
	if !exists {
		return nil, false
	}
	
	// Validate type compatibility
	if cfg.Type != "" && registration.Type != "" && cfg.Type != registration.Type {
		return nil, false
	}
	
	// Create the provider using the registered constructor
	client := registration.Constructor(opts)
	return client, true
}

// MustRegisterProvider registers a provider and panics on error (for init() functions)
func MustRegisterProvider(registration *ProviderRegistration) {
	if err := RegisterProvider(registration); err != nil {
		panic(fmt.Sprintf("failed to register provider %s: %v", registration.ID, err))
	}
}


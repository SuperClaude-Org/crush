package config

import (
	"github.com/charmbracelet/catwalk/pkg/catwalk"
	"github.com/charmbracelet/crush/internal/auth"
)

// OAuthProvider represents a provider that uses OAuth authentication
type OAuthProvider struct {
	ID       string
	Name     string
	Type     catwalk.Type
	Models   []catwalk.Model
	BaseURL  string
	AuthManager *auth.AuthManager
}

// OAuthProviderRegistry manages OAuth providers
type OAuthProviderRegistry struct {
	providers map[string]OAuthProvider
}

var oauthRegistry = &OAuthProviderRegistry{
	providers: make(map[string]OAuthProvider),
}

// RegisterOAuthProvider registers a new OAuth provider
func RegisterOAuthProvider(provider OAuthProvider) {
	oauthRegistry.providers[provider.ID] = provider
}

// GetOAuthProviders returns all registered OAuth providers
func GetOAuthProviders(dataDirectory string) []OAuthProvider {
	authManager := auth.NewAuthManager(dataDirectory)
	
	var providers []OAuthProvider
	for _, provider := range oauthRegistry.providers {
		// Create a copy with the auth manager
		providerCopy := provider
		providerCopy.AuthManager = authManager
		providers = append(providers, providerCopy)
	}
	
	return providers
}

// GetOAuthProvider returns a specific OAuth provider by ID
func GetOAuthProvider(id string, dataDirectory string) (*OAuthProvider, bool) {
	provider, exists := oauthRegistry.providers[id]
	if !exists {
		return nil, false
	}
	
	// Create a copy with the auth manager
	providerCopy := provider
	providerCopy.AuthManager = auth.NewAuthManager(dataDirectory)
	
	return &providerCopy, true
}

// HasOAuthCredentials checks if an OAuth provider has valid credentials
func (p *OAuthProvider) HasOAuthCredentials() bool {
	if p.AuthManager == nil {
		return false
	}
	
	switch p.ID {
	case "claudesub":
		return p.AuthManager.HasClaudeSubAuth()
	default:
		return false
	}
}

// ToDisplayProvider converts an OAuth provider to a catwalk.Provider for TUI display
func (p *OAuthProvider) ToDisplayProvider() catwalk.Provider {
	name := p.Name
	if !p.HasOAuthCredentials() {
		name += " (Auth Required)"
	}
	
	// Convert models to OAuth models with zero cost
	oauthModels := make([]catwalk.Model, len(p.Models))
	for i, model := range p.Models {
		oauthModels[i] = catwalk.Model{
			ID:                     model.ID,
			Name:                   model.Name,
			CostPer1MIn:            0, // OAuth providers are subscription-based
			CostPer1MOut:           0,
			CostPer1MInCached:      0,
			CostPer1MOutCached:     0,
			ContextWindow:          model.ContextWindow,
			DefaultMaxTokens:       model.DefaultMaxTokens,
			CanReason:              model.CanReason,
			HasReasoningEffort:     model.HasReasoningEffort,
			DefaultReasoningEffort: model.DefaultReasoningEffort,
			SupportsImages:         model.SupportsImages,
		}
	}
	
	return catwalk.Provider{
		ID:     catwalk.InferenceProvider(p.ID),
		Name:   name,
		Type:   p.Type,
		Models: oauthModels,
	}
}

// RegisterBuiltinOAuthProviders registers the built-in OAuth providers
func RegisterBuiltinOAuthProviders() {
	// Register Claude Subscription provider
	RegisterOAuthProvider(OAuthProvider{
		ID:      "claudesub",
		Name:    "Claude Max/Pro Subscription",
		Type:    catwalk.TypeAnthropic,
		BaseURL: "https://api.anthropic.com/v1",
		Models: []catwalk.Model{
			{
				ID:                     "claude-opus-4-1-20250805",
				Name:                   "Claude Opus 4.1",
				ContextWindow:          200000,
				DefaultMaxTokens:       32000,
				CanReason:              true,
				HasReasoningEffort:     false,
				DefaultReasoningEffort: "",
				SupportsImages:         true,
			},
			{
				ID:                     "claude-opus-4-20250514",
				Name:                   "Claude Opus 4",
				ContextWindow:          200000,
				DefaultMaxTokens:       32000,
				CanReason:              true,
				HasReasoningEffort:     false,
				DefaultReasoningEffort: "",
				SupportsImages:         true,
			},
			{
				ID:                     "claude-sonnet-4-20250514",
				Name:                   "Claude Sonnet 4",
				ContextWindow:          200000,
				DefaultMaxTokens:       8192,
				CanReason:              true,
				HasReasoningEffort:     true,
				DefaultReasoningEffort: "medium",
				SupportsImages:         true,
			},
			{
				ID:                     "claude-3-7-sonnet-20250219",
				Name:                   "Claude 3.7 Sonnet",
				ContextWindow:          200000,
				DefaultMaxTokens:       8192,
				CanReason:              true,
				HasReasoningEffort:     true,
				DefaultReasoningEffort: "medium",
				SupportsImages:         true,
			},
			{
				ID:                     "claude-3-5-sonnet-20241022",
				Name:                   "Claude 3.5 Sonnet (New)",
				ContextWindow:          200000,
				DefaultMaxTokens:       8192,
				CanReason:              true,
				HasReasoningEffort:     true,
				DefaultReasoningEffort: "medium",
				SupportsImages:         true,
			},
			{
				ID:                     "claude-3-5-sonnet-20240620",
				Name:                   "Claude 3.5 Sonnet (Old)",
				ContextWindow:          200000,
				DefaultMaxTokens:       8192,
				CanReason:              true,
				HasReasoningEffort:     true,
				DefaultReasoningEffort: "medium",
				SupportsImages:         true,
			},
			{
				ID:               "claude-3-5-haiku-20241022",
				Name:             "Claude 3.5 Haiku",
				ContextWindow:    200000,
				DefaultMaxTokens: 5000,
				CanReason:        false,
				SupportsImages:   true,
			},
		},
	})
}

// Initialize OAuth providers
func init() {
	RegisterBuiltinOAuthProviders()
}
package provider

import (
	"context"

	"github.com/charmbracelet/catwalk/pkg/catwalk"
	"github.com/charmbracelet/crush/internal/auth"
	"github.com/charmbracelet/crush/internal/config"
)

// OAuthProvider defines the interface for providers that support OAuth authentication
type OAuthProvider interface {
	// SupportsOAuth returns true if the provider supports OAuth authentication
	SupportsOAuth() bool
	
	// HasOAuthCredentials returns true if OAuth credentials are available
	HasOAuthCredentials() bool
	
	// GetOAuthModels returns models with OAuth-specific configuration (e.g., zero cost)
	GetOAuthModels() []catwalk.Model
	
	// RequiresOAuthSetup returns true if this provider requires OAuth to function
	RequiresOAuthSetup() bool
	
	// GetAuthManager returns the auth manager for this provider
	GetAuthManager() *auth.AuthManager
}

// OAuthCapableProvider is a helper interface that extends ProviderClient with OAuth capabilities
type OAuthCapableProvider interface {
	ProviderClient
	OAuthProvider
}

// IsOAuthProvider checks if a provider client implements OAuth support
func IsOAuthProvider(client ProviderClient) (OAuthProvider, bool) {
	oauthProvider, ok := client.(OAuthProvider)
	return oauthProvider, ok
}

// GetOAuthProviderFromConfig creates an OAuth provider from configuration if supported
func GetOAuthProviderFromConfig(cfg config.ProviderConfig, client ProviderClient) (OAuthProvider, bool) {
	if oauthProvider, ok := IsOAuthProvider(client); ok {
		return oauthProvider, true
	}
	return nil, false
}

// HasValidOAuthCredentials checks if a provider has valid OAuth credentials
func HasValidOAuthCredentials(ctx context.Context, client ProviderClient) bool {
	oauthProvider, ok := IsOAuthProvider(client)
	if !ok {
		return false
	}
	
	if !oauthProvider.SupportsOAuth() {
		return false
	}
	
	if !oauthProvider.HasOAuthCredentials() {
		return false
	}
	
	// Could add token validation here if needed
	return true
}

// RequiresOAuthAuthentication checks if a provider requires OAuth to function
func RequiresOAuthAuthentication(client ProviderClient) bool {
	oauthProvider, ok := IsOAuthProvider(client)
	if !ok {
		return false
	}
	
	return oauthProvider.RequiresOAuthSetup()
}
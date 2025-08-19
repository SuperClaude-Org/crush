package config

import "github.com/charmbracelet/catwalk/pkg/catwalk"

// IsOAuthProvider checks if a provider ID supports OAuth authentication
// This is a simple, direct approach without unnecessary abstractions
func IsOAuthProvider(providerID string) bool {
	switch providerID {
	case "claudesub":
		return true
	default:
		return false
	}
}

// ListOAuthProviders returns all OAuth-capable provider IDs
func ListOAuthProviders() []string {
	return []string{"claudesub"}
}

// GetDefaultOAuthModels returns default models for OAuth providers
// This provides fallback models when no config models are specified
func GetDefaultOAuthModels(providerID string) []catwalk.Model {
	switch providerID {
	case "claudesub":
		return []catwalk.Model{
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
		}
	default:
		return nil
	}
}
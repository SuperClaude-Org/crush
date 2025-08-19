package provider

import (
	"context"
	"log/slog"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/charmbracelet/catwalk/pkg/catwalk"
	"github.com/charmbracelet/crush/internal/auth"
	"github.com/charmbracelet/crush/internal/config"
	"github.com/charmbracelet/crush/internal/llm/tools"
	"github.com/charmbracelet/crush/internal/message"
)

const (
	// systemPromptPrefix is required for OAuth authentication to work
	// This identifies the client as Claude Code to the Anthropic API
	systemPromptPrefix = "You are Claude Code, Anthropic's official CLI for Claude."
)

type claudeSubClient struct {
	providerOptions providerClientOptions
	anthropicClient AnthropicClient
	authManager     *auth.AuthManager
	useOAuth        bool
}

// claudeSubClient implements ProviderClient interface for OAuth-enabled Claude provider

func newClaudeSubClient(opts providerClientOptions) ProviderClient {
	cfg := config.Get()
	if cfg == nil {
		slog.Error("Configuration not loaded for claudesub provider")
		return nil
	}
	
	// Automatically inject system prompt prefix for OAuth authentication
	opts.systemPromptPrefix = systemPromptPrefix
	
	authManager := auth.NewAuthManager(cfg.Options.DataDirectory)
	
	useOAuth := authManager.HasClaudeSubAuth()
	
	var anthropicClient AnthropicClient
	
	if useOAuth {
		slog.Info("Using OAuth authentication for claudesub provider")
		httpClient, err := createClaudeSubHTTPClient()
		if err != nil {
			slog.Error("Failed to create OAuth HTTP client", "error", err)
			useOAuth = false
			anthropicClient = newAnthropicClient(opts, AnthropicClientTypeNormal)
		} else {
			anthropicOpts := []option.RequestOption{
				option.WithHTTPClient(httpClient),
			}
			
			// Only set custom base URL if it's not the default Anthropic URL
			// The Anthropic SDK already handles the default URL correctly
			if opts.baseURL != "" && opts.baseURL != "https://api.anthropic.com/v1" {
				// Strip /v1 suffix if present to avoid doubling
				baseURL := opts.baseURL
				if strings.HasSuffix(baseURL, "/v1") {
					baseURL = strings.TrimSuffix(baseURL, "/v1")
				}
				anthropicOpts = append(anthropicOpts, option.WithBaseURL(baseURL))
			}
			
			// Add extra headers
			for key, header := range opts.extraHeaders {
				anthropicOpts = append(anthropicOpts, option.WithHeaderAdd(key, header))
			}
			
			// Add extra body parameters
			for key, value := range opts.extraBody {
				anthropicOpts = append(anthropicOpts, option.WithJSONSet(key, value))
			}
			
			client := newAnthropicClientWithOptions(opts, anthropicOpts)
			anthropicClient = client
		}
	} else {
		slog.Info("OAuth not available, using API key authentication for claudesub provider")
		anthropicClient = newAnthropicClient(opts, AnthropicClientTypeNormal)
	}
	
	return &claudeSubClient{
		providerOptions: opts,
		anthropicClient: anthropicClient,
		authManager:     authManager,
		useOAuth:        useOAuth,
	}
}

// newAnthropicClientWithOptions creates an Anthropic client with custom options
func newAnthropicClientWithOptions(opts providerClientOptions, anthropicOpts []option.RequestOption) AnthropicClient {
	return &anthropicClient{
		providerOptions: opts,
		tp:              AnthropicClientTypeNormal,
		client:          anthropic.NewClient(anthropicOpts...),
	}
}


func (c *claudeSubClient) send(ctx context.Context, messages []message.Message, tools []tools.BaseTool) (*ProviderResponse, error) {
	return c.anthropicClient.send(ctx, messages, tools)
}

func (c *claudeSubClient) stream(ctx context.Context, messages []message.Message, tools []tools.BaseTool) <-chan ProviderEvent {
	return c.anthropicClient.stream(ctx, messages, tools)
}

func (c *claudeSubClient) Model() catwalk.Model {
	// Return the selected model or fall back to the configured model
	if c.providerOptions.modelType != "" {
		return c.providerOptions.model(c.providerOptions.modelType)
	}
	
	// Return underlying model if no specific model type is set
	return c.anthropicClient.Model()
}

// OAuth Provider Interface Implementation
func (c *claudeSubClient) SupportsOAuth() bool {
	return true
}

func (c *claudeSubClient) HasOAuthCredentials() bool {
	return c.authManager.HasClaudeSubAuth()
}

func (c *claudeSubClient) GetOAuthModels() []catwalk.Model {
	if !c.HasOAuthCredentials() {
		return nil
	}
	
	// If models are configured, use them; otherwise use defaults
	var sourceModels []catwalk.Model
	if len(c.providerOptions.config.Models) > 0 {
		// Convert provider models to OAuth models with zero cost
		sourceModels = c.providerOptions.config.Models
	} else {
		// Use default model set for claudesub OAuth provider
		sourceModels = getDefaultClaudeSubModels()
	}
	
	oauthModels := make([]catwalk.Model, len(sourceModels))
	for i, model := range sourceModels {
		oauthModels[i] = catwalk.Model{
			ID:                     model.ID,
			Name:                   model.Name,
			CostPer1MIn:            0,
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
	
	return oauthModels
}

func (c *claudeSubClient) RequiresOAuthSetup() bool {
	return true
}

func (c *claudeSubClient) GetAuthManager() *auth.AuthManager {
	return c.authManager
}

// getDefaultClaudeSubModels returns the default model set for claudesub OAuth provider
func getDefaultClaudeSubModels() []catwalk.Model {
	return []catwalk.Model{
		{
			ID:                     "claude-opus-4-1-20250805",
			Name:                   "Claude Opus 4.1",
			ContextWindow:          200000,
			DefaultMaxTokens:       32000,
			CanReason:              true,
			HasReasoningEffort:     true,
			DefaultReasoningEffort: "",
			SupportsImages:         true,
		},
		{
			ID:                     "claude-opus-4-20250514",
			Name:                   "Claude Opus 4",
			ContextWindow:          200000,
			DefaultMaxTokens:       32000,
			CanReason:              true,
			HasReasoningEffort:     true,
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
}

// Register claudesub provider at package initialization
func init() {
	MustRegisterProvider(&ProviderRegistration{
		ID:            "claudesub",
		Name:          "Claude Max/Pro Subscription",
		Type:          catwalk.TypeAnthropic,
		SupportsOAuth: true,
		Constructor: func(opts providerClientOptions) ProviderClient {
			return newClaudeSubClient(opts)
		},
	})
}


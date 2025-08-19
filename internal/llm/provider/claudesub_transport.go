package provider

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/charmbracelet/crush/internal/auth"
	"github.com/charmbracelet/crush/internal/config"
)


type claudeSubTransport struct {
	base        http.RoundTripper
	authManager *auth.AuthManager
}

func newClaudeSubTransport(authManager *auth.AuthManager) http.RoundTripper {
	return &claudeSubTransport{
		base:        http.DefaultTransport,
		authManager: authManager,
	}
}

func (t *claudeSubTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	
	// Clone the request to avoid modifying the original
	clonedReq := req.Clone(req.Context())
	
	if err := t.modifyRequestHeaders(clonedReq); err != nil {
		return nil, fmt.Errorf("failed to modify request headers: %w", err)
	}
	
	// Make the request
	resp, err := t.base.RoundTrip(clonedReq)
	if err != nil {
		return nil, err
	}
	
	// Handle 401 Unauthorized - attempt token refresh
	if resp.StatusCode == http.StatusUnauthorized {
		// Close the failed response
		resp.Body.Close()
		
		// Try to refresh the token
		if refreshErr := t.authManager.RefreshClaudeSubToken(req.Context()); refreshErr != nil {
			slog.Error("Failed to refresh token", "error", refreshErr)
			return nil, fmt.Errorf("authentication failed and token refresh failed: %w", refreshErr)
		}
		
		// Retry the request with refreshed token
		clonedReq2 := req.Clone(req.Context())
		
		if err := t.modifyRequestHeaders(clonedReq2); err != nil {
			return nil, fmt.Errorf("failed to modify request headers after refresh: %w", err)
		}
		
		resp, err = t.base.RoundTrip(clonedReq2)
		if err != nil {
			return nil, err
		}
	}
	
	return resp, nil
}

func (t *claudeSubTransport) modifyRequestHeaders(req *http.Request) error {
	// Get valid access token (will refresh if needed)
	accessToken, err := t.authManager.GetValidAccessToken(req.Context())
	if err != nil {
		return fmt.Errorf("failed to get valid access token: %w", err)
	}
	
	// Remove any existing API key headers
	req.Header.Del("x-api-key")
	req.Header.Del("X-Api-Key")
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	
	req.Header.Del("Anthropic-Beta") 
	req.Header.Del("anthropic-beta")
	req.Header.Add("anthropic-beta", auth.BetaHeaders)
	req.Header.Del("User-Agent")
	req.Header.Del("X-Stainless-Os")
	req.Header.Del("X-Stainless-Lang")
	req.Header.Del("X-Stainless-Retry-Count")
	req.Header.Del("X-Stainless-Arch")
	req.Header.Del("X-Stainless-Runtime")
	req.Header.Del("X-Stainless-Runtime-Version")
	req.Header.Del("X-Stainless-Package-Version")
	
	return nil
}

func createClaudeSubHTTPClient() (*http.Client, error) {
	cfg := config.Get()
	if cfg == nil {
		return nil, fmt.Errorf("configuration not loaded")
	}
	
	authManager := auth.NewAuthManager(cfg.Options.DataDirectory)
	
	if !authManager.HasClaudeSubAuth() {
		return nil, fmt.Errorf("no OAuth credentials found - run 'crush auth login' first")
	}
	
	transport := newClaudeSubTransport(authManager)
	
	// Add debug logging if enabled
	if cfg.Options.Debug {
		transport = &HTTPDebugTransport{
			Transport: transport,
		}
	}
	
	return &http.Client{
		Transport: transport,
	}, nil
}

// HTTPDebugTransport wraps another transport with debug logging
type HTTPDebugTransport struct {
	Transport http.RoundTripper
}

func (t *HTTPDebugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.Transport.RoundTrip(req)
}

// filterSensitiveHeaders removes sensitive headers from debug output
func filterSensitiveHeaders(headers http.Header) http.Header {
	filtered := make(http.Header)
	for k, v := range headers {
		switch strings.ToLower(k) {
		case "authorization", "x-api-key":
			filtered[k] = []string{"[REDACTED]"}
		default:
			filtered[k] = v
		}
	}
	return filtered
}
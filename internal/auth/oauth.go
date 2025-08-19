package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	ClientID     = "9d1c250a-e61b-44d9-88ed-5944d1962f5e"
	RedirectURI  = "https://console.anthropic.com/oauth/code/callback"
	Scope        = "org:create_api_key user:profile user:inference"
	AuthURL      = "https://claude.ai/oauth/authorize"
	TokenURL     = "https://console.anthropic.com/v1/oauth/token"
	BetaHeaders  = "oauth-2025-04-20,claude-code-20250219,interleaved-thinking-2025-05-14,fine-grained-tool-streaming-2025-05-14"
)

type OAuthFlow struct {
	httpClient *http.Client
	pkce       *PKCEChallenge
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope,omitempty"` // Granted scopes
}

func NewOAuthFlow() *OAuthFlow {
	return &OAuthFlow{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (o *OAuthFlow) GenerateAuthURL() (string, error) {
	pkce, err := GeneratePKCEChallenge()
	if err != nil {
		return "", fmt.Errorf("failed to generate PKCE challenge: %w", err)
	}
	o.pkce = pkce

	params := url.Values{}
	params.Set("client_id", ClientID)
	params.Set("response_type", "code")
	params.Set("redirect_uri", RedirectURI)
	params.Set("scope", Scope)
	params.Set("code_challenge", pkce.Challenge)
	params.Set("code_challenge_method", pkce.Method)
	params.Set("code", "true")
	params.Set("state", pkce.Verifier) // Use verifier as state

	authURL := AuthURL + "?" + params.Encode()
	return authURL, nil
}

func (o *OAuthFlow) OpenBrowser(authURL string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", authURL}
	case "darwin":
		cmd = "open"
		args = []string{authURL}
	case "linux":
		cmd = "xdg-open"
		args = []string{authURL}
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return exec.Command(cmd, args...).Start()
}

func (o *OAuthFlow) ExchangeCodeForTokens(ctx context.Context, codeWithState string) (*TokenResponse, error) {
	if o.pkce == nil {
		return nil, fmt.Errorf("PKCE challenge not initialized - call GenerateAuthURL first")
	}

	// Parse the code#state format
	parts := strings.Split(codeWithState, "#")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid code format - expected 'code#state', got: %s", codeWithState)
	}

	authCode := parts[0]
	state := parts[1]

	// Verify state matches our verifier
	if state != o.pkce.Verifier {
		return nil, fmt.Errorf("state mismatch - possible CSRF attack")
	}

	requestData := map[string]interface{}{
		"grant_type":    "authorization_code",
		"client_id":     ClientID,
		"code":          authCode,
		"redirect_uri":  RedirectURI,
		"code_verifier": o.pkce.Verifier,
		"state":         state,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal token request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", TokenURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token exchange request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}

func (o *OAuthFlow) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	requestData := map[string]interface{}{
		"grant_type":    "refresh_token",
		"client_id":     ClientID,
		"refresh_token": refreshToken,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal refresh request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", TokenURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read refresh response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token refresh failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse refresh response: %w", err)
	}

	return &tokenResp, nil
}
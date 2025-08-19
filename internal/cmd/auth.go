package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/crush/internal/auth"
	"github.com/charmbracelet/crush/internal/config"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication for AI providers",
	Long: `Manage authentication credentials for AI providers.
Supports OAuth authentication for Claude Pro/Max subscriptions and API key management.`,
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with an AI provider",
	Long: `Start the authentication process with an AI provider.
Currently supports Claude Pro/Max OAuth authentication for subscription access.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAuthLogin(cmd.Context())
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored authentication credentials",
	Long:  `Remove stored authentication credentials for AI providers.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAuthLogout(cmd.Context())
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	Long:  `Display the current authentication status for all configured providers.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAuthStatus(cmd.Context())
	},
}

func init() {
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authLogoutCmd)
	authCmd.AddCommand(authStatusCmd)
	
	rootCmd.AddCommand(authCmd)
}

// runAuthLogin handles the authentication login process
func runAuthLogin(ctx context.Context) error {
	cfg, err := config.Load("", false)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	authManager := auth.NewAuthManager(cfg.Options.DataDirectory)
	
	// Check if already authenticated
	if authManager.HasClaudeSubAuth() {
		valid, err := authManager.IsClaudeSubTokenValid()
		if err != nil {
			slog.Warn("Error checking token validity", "error", err)
		} else if valid {
			fmt.Println("✅ Already authenticated with Claude Pro/Max subscription")
			fmt.Println("💡 Use 'crush auth logout' to sign out and authenticate with a different account")
			return nil
		}
	}

	fmt.Println("🔐 Claude Pro/Max Subscription Authentication")
	fmt.Println("📝 This will authenticate you with your Claude Pro or Max subscription")
	fmt.Println("")
	
	// Show provider selection menu
	fmt.Println("Available authentication methods:")
	fmt.Println("  1. Claude Pro/Max (OAuth via claude.ai) [RECOMMENDED]")
	fmt.Println("")
	
	// For now, we only support Claude Pro/Max, so auto-select it
	fmt.Printf("🚀 Starting OAuth authentication with Claude Pro/Max...\n\n")
	
	return authenticateClaudeSub(ctx, authManager)
}

func authenticateClaudeSub(ctx context.Context, authManager *auth.AuthManager) error {
	oauthFlow := auth.NewOAuthFlow()
	
	// Generate authorization URL
	authURL, err := oauthFlow.GenerateAuthURL()
	if err != nil {
		return fmt.Errorf("failed to generate authorization URL: %w", err)
	}
	
	fmt.Println("📱 Opening your browser for authentication...")
	fmt.Printf("🌐 Auth URL: %s\n\n", authURL)
	
	// Try to open browser
	if err := oauthFlow.OpenBrowser(authURL); err != nil {
		slog.Debug("Failed to open browser", "error", err)
		fmt.Println("⚠️  Could not open browser automatically")
		fmt.Printf("🔗 Please manually open this URL in your browser:\n%s\n\n", authURL)
	} else {
		fmt.Println("✅ Opened browser for authentication")
		fmt.Println("")
	}
	
	// Instructions for user
	fmt.Println("📋 Instructions:")
	fmt.Println("   1. Complete the OAuth authorization in your browser")
	fmt.Println("   2. You'll be redirected to a callback page")
	fmt.Println("   3. Copy the authorization code from the callback page")
	fmt.Println("   4. The code will be in format: code#state")
	fmt.Println("")
	
	// Prompt for authorization code
	fmt.Print("📝 Paste the authorization code here: ")
	reader := bufio.NewReader(os.Stdin)
	codeInput, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read authorization code: %w", err)
	}
	
	codeInput = strings.TrimSpace(codeInput)
	if codeInput == "" {
		return fmt.Errorf("no authorization code provided")
	}
	
	fmt.Println("\n🔄 Exchanging authorization code for tokens...")
	
	// Exchange code for tokens
	tokenResp, err := oauthFlow.ExchangeCodeForTokens(ctx, codeInput)
	if err != nil {
		return fmt.Errorf("token exchange failed: %w", err)
	}
	
	fmt.Println("💾 Storing authentication credentials...")
	
	// Store credentials
	if err := authManager.StoreClaudeSubCredentials(tokenResp); err != nil {
		return fmt.Errorf("failed to store credentials: %w", err)
	}
	
	fmt.Println("✅ Successfully authenticated with Claude Pro/Max!")
	fmt.Printf("⏰ Token expires: %s\n", time.Now().Add(time.Duration(tokenResp.ExpiresIn)*time.Second).Format("2006-01-02 15:04:05"))
	fmt.Println("🎉 You can now use the claudesub provider with your subscription")
	fmt.Println("")
	fmt.Println("💡 Configure your models to use the 'claudesub' provider:")
	fmt.Printf("   crush config set models.large.provider claudesub\n")
	
	return nil
}

// runAuthLogout handles the authentication logout process
func runAuthLogout(ctx context.Context) error {
	cfg, err := config.Load("", false)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	authManager := auth.NewAuthManager(cfg.Options.DataDirectory)
	
	// Check if authenticated
	if !authManager.HasClaudeSubAuth() {
		fmt.Println("ℹ️  No Claude Pro/Max authentication found")
		return nil
	}
	
	fmt.Println("🔓 Signing out of Claude Pro/Max subscription...")
	
	// Clear credentials
	if err := authManager.ClearClaudeSubCredentials(); err != nil {
		return fmt.Errorf("failed to clear credentials: %w", err)
	}
	
	fmt.Println("✅ Successfully signed out")
	fmt.Println("💡 Use 'crush auth login' to authenticate again")
	
	return nil
}

// runAuthStatus displays the current authentication status
func runAuthStatus(ctx context.Context) error {
	cfg, err := config.Load("", false)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	authManager := auth.NewAuthManager(cfg.Options.DataDirectory)
	
	fmt.Println("🔐 Authentication Status")
	fmt.Println("========================")
	fmt.Println("")
	
	// Check Claude Pro/Max status
	if authManager.HasClaudeSubAuth() {
		creds, err := authManager.GetClaudeSubCredentials()
		if err != nil {
			fmt.Printf("❌ Claude Pro/Max: Error reading credentials (%v)\n", err)
		} else {
			valid, err := authManager.IsClaudeSubTokenValid()
			if err != nil {
				fmt.Printf("⚠️  Claude Pro/Max: Cannot verify token validity (%v)\n", err)
			} else if valid {
				expiresAt := time.UnixMilli(creds.Expires)
				fmt.Printf("✅ Claude Pro/Max: Authenticated (expires %s)\n", expiresAt.Format("2006-01-02 15:04:05"))
			} else {
				fmt.Println("⚠️  Claude Pro/Max: Token expired (will auto-refresh on next use)")
			}
		}
	} else {
		fmt.Println("❌ Claude Pro/Max: Not authenticated")
	}
	
	fmt.Println("")
	fmt.Println("💡 Use 'crush auth login' to authenticate")
	fmt.Println("💡 Use 'crush auth logout' to sign out")
	
	return nil
}
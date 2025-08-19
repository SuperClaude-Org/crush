# Claude Subscription Provider TUI Display Fix

## Issue
The Claude Subscription (claudesub) OAuth provider was not appearing in the TUI provider list despite being properly registered.

## Root Cause
In `internal/config/load.go`, the `configureProviders()` function was incorrectly deleting OAuth providers during validation. Specifically:

- Lines 273-277 required all custom providers to have a `BaseURL`
- OAuth providers like `claudesub` don't need a `BaseURL` since they authenticate via OAuth tokens
- This caused the provider to be deleted before the TUI could display it

## Solution
Modified the validation logic in `internal/config/load.go:274` to exempt OAuth providers from the BaseURL requirement:

```go
// OAuth providers don't need a BaseURL as they use OAuth authentication
if providerConfig.BaseURL == "" && !IsOAuthProvider(id) {
    slog.Warn("Skipping custom provider due to missing API endpoint", "provider", id)
    c.Providers.Del(id)
    continue
}
```

## Verification
After the fix:
- No more warnings about "Skipping custom provider due to missing API endpoint" for claudesub
- Logs show "Using OAuth authentication for claudesub provider" indicating proper recognition
- Auth status command works correctly: `./crush auth status`
- Provider properly loads and validates OAuth credentials

## Files Modified
- `internal/config/load.go` (lines 273-278): Added OAuth provider exemption for BaseURL requirement

## Testing
```bash
# Build with fix
go build -o ./crush .

# Verify authentication works
./crush auth status

# Check logs for proper provider loading
./crush logs --tail 30
```

The claudesub provider should now appear in the TUI provider list with appropriate authentication status indicators.
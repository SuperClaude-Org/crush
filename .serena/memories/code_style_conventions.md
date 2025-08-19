# Crush Code Style and Conventions

## Go Style Guidelines
- Follows standard Go conventions and best practices
- Uses `gofumpt` for formatting (stricter than gofmt)
- Package names are lowercase, single word when possible
- Exported functions use PascalCase, unexported use camelCase

## Code Organization
- Internal packages under `internal/` directory
- Clear separation of concerns by package
- Commands in `internal/cmd/`
- Business logic in `internal/app/`
- Provider implementations in `internal/llm/provider/`

## Naming Conventions
- **Files**: snake_case.go (e.g., `claudesub.go`, `oauth_provider.go`)
- **Types**: PascalCase (e.g., `ClaudeSubClient`, `AuthManager`)
- **Functions**: PascalCase for exported, camelCase for unexported
- **Constants**: PascalCase or UPPER_CASE for package-level constants
- **Variables**: camelCase (e.g., `authManager`, `tokenResp`)

## Error Handling
- Return errors as last return value
- Use `fmt.Errorf` for error wrapping with context
- Log errors at appropriate levels with structured logging
- Use specific error types when needed

## Logging
- Uses structured logging with `log/slog`
- Debug level for detailed debugging information
- Info level for normal operation events
- Warn level for unexpected but recoverable conditions
- Error level for errors that need attention

## Configuration
- JSON-based configuration with schema validation
- Environment variable support with `godotenv`
- Hierarchical config loading (project > user > global)

## Testing
- Test files end with `_test.go`
- Use `testify` for assertions
- Table-driven tests where appropriate
- Mock interfaces for testing

## Documentation
- Package documentation in doc.go files
- Exported functions and types have godoc comments
- README.md for project overview
- Examples in godoc comments when helpful

## Dependencies
- Minimal external dependencies
- Prefer standard library when possible
- Use established, well-maintained packages
- Pin to specific versions in go.mod

## CLI Design
- Uses Cobra framework for command structure
- Consistent flag naming across commands
- Help text and examples for all commands
- Non-zero exit codes for errors
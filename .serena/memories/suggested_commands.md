# Crush Development Commands

## Building and Testing
```bash
# Build the application
go build -o ./crush .

# Run tests
task test
# or
go test ./...

# Install application
task install
# or 
go install -v .
```

## Code Quality
```bash
# Install linter
task lint:install

# Run linters
task lint

# Run linters with fixes
task lint-fix

# Format code
task fmt
# or
gofumpt -w .
```

## Development and Debugging
```bash
# Run with profiling enabled
task dev
# or
CRUSH_PROFILE=true go run .

# Run with debug logging
./crush --debug
./crush run "prompt" --debug

# CPU profiling (after starting with profiling)
task profile:cpu

# Heap profiling
task profile:heap

# Allocations profiling  
task profile:allocs
```

## Application Usage
```bash
# Interactive mode
./crush

# Non-interactive mode
./crush run "your prompt here"

# Quiet mode (no spinner)
./crush run "prompt" --quiet

# Change working directory
./crush --cwd /path/to/project

# Auto-accept permissions (dangerous)
./crush --yolo
```

## Authentication (New Claude Subscription Feature)
```bash
# Check authentication status
./crush auth status

# Authenticate with Claude Pro/Max
./crush auth login

# Sign out
./crush auth logout
```

## Logging
```bash
# View recent logs
./crush logs

# Follow logs in real-time
./crush logs --follow

# Show last N lines
./crush logs --tail 500
```

## Schema Generation
```bash
# Generate JSON schema for configuration
task schema
# or
go run main.go schema > schema.json
```

## Shell Completion
```bash
# Generate bash completion
./crush completion bash > /etc/bash_completion.d/crush

# Generate for other shells
./crush completion [bash|zsh|fish|powershell]
```

## System Utilities (Linux)
Standard Linux commands are available:
- `git` - Version control
- `ls`, `cd`, `grep`, `find` - File operations  
- `cat`, `less`, `tail` - File viewing
- `ps`, `top`, `htop` - Process management
- `curl`, `wget` - HTTP requests
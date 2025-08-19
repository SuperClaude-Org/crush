# Crush Commands and Flags Reference

## Main Commands

### Root Command (Interactive Mode)
```bash
crush [--flags]
```
**Flags:**
- `-c, --cwd <path>` - Current working directory
- `-d, --debug` - Enable debug logging
- `-h, --help` - Show help
- `-v, --version` - Show version
- `-y, --yolo` - Auto-accept all permissions (dangerous mode)

### Run Command (Non-Interactive)
```bash
crush run [prompt...] [--flags]
```
**Purpose**: Execute a single prompt and exit
**Flags:**
- `-c, --cwd <path>` - Current working directory  
- `-d, --debug` - Enable debug logging
- `-h, --help` - Show help
- `-q, --quiet` - Hide spinner

**Examples:**
```bash
crush run "Explain this code"
echo "What does this do?" | crush run
crush run -q "Generate README"
```

### Auth Commands
```bash
crush auth [command] [--flags]
```

#### auth login
```bash
crush auth login [--flags]
```
**Purpose**: Authenticate with Claude Pro/Max subscription via OAuth
**Flags:**
- `-c, --cwd <path>` - Current working directory
- `-d, --debug` - Enable debug logging  
- `-h, --help` - Show help

#### auth logout
```bash
crush auth logout [--flags]
```
**Purpose**: Remove stored authentication credentials
**Flags:**
- `-c, --cwd <path>` - Current working directory
- `-d, --debug` - Enable debug logging
- `-h, --help` - Show help

#### auth status
```bash
crush auth status [--flags]
```
**Purpose**: Show current authentication status
**Flags:**
- `-c, --cwd <path>` - Current working directory
- `-d, --debug` - Enable debug logging
- `-h, --help` - Show help

### Logs Command
```bash
crush logs [--flags]
```
**Purpose**: View crush logs for debugging and monitoring
**Flags:**
- `-c, --cwd <path>` - Current working directory
- `-d, --debug` - Enable debug logging
- `-f, --follow` - Follow log output in real-time
- `-h, --help` - Show help
- `-t, --tail <number>` - Show only last N lines (default: 1000)

**Examples:**
```bash
crush logs                  # Last 1000 lines
crush logs --tail 500      # Last 500 lines  
crush logs --follow        # Real-time following
```

### Schema Command
```bash
crush schema [--flags]
```
**Purpose**: Generate JSON schema for configuration
**Flags:**
- `-c, --cwd <path>` - Current working directory
- `-d, --debug` - Enable debug logging
- `-h, --help` - Show help

### Completion Command
```bash
crush completion [shell] [--flags]
```
**Purpose**: Generate shell completion scripts
**Shells**: bash, zsh, fish, powershell
**Flags:**
- `-c, --cwd <path>` - Current working directory
- `-d, --debug` - Enable debug logging  
- `-h, --help` - Show help
- `--no-descriptions` - Disable completion descriptions (bash only)

**Examples:**
```bash
crush completion bash > /etc/bash_completion.d/crush
crush completion zsh > ~/.zsh/completions/_crush
```

### Help Command
```bash
crush help [command]
```
**Purpose**: Get help about any command

## Global Flags (Available on All Commands)
- `-c, --cwd <path>` - Set current working directory
- `-d, --debug` - Enable debug logging
- `-h, --help` - Show help for command

## Environment Variables
- `CRUSH_PROFILE=true` - Enable profiling (pprof on localhost:6060)
- `CLAUDE_SUB_API_KEY` - API key for testing (development use)
- Various provider API keys (see README.md for full list)

## Exit Codes
- `0` - Success
- `1` - General error
- Other non-zero codes for specific error conditions

## Special Modes
- **Interactive Mode**: Default when no command specified
- **Non-Interactive Mode**: Using `crush run` command  
- **Debug Mode**: Enabled with `--debug` flag
- **YOLO Mode**: Auto-accept permissions with `--yolo` flag
- **Quiet Mode**: Suppress spinner with `--quiet` (run command only)
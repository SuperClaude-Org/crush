# Crush Project Overview

## Purpose
Crush is a powerful terminal-based AI assistant for software development. It's a fork that adds Claude Subscription (Pro/Max) OAuth authentication support through a new `claudesub` provider.

Key features:
- Multi-model support (OpenAI, Anthropic, Groq, OpenRouter, etc.)
- Session-based conversations with context preservation
- LSP integration for code analysis
- Model Context Protocol (MCP) server support
- OAuth authentication for Claude Pro/Max subscriptions
- Cross-platform terminal application

## Tech Stack
- **Language**: Go 1.25.0
- **Framework**: Cobra CLI with Bubble Tea TUI
- **Database**: SQLite (go-sqlite3)
- **Authentication**: OAuth 2.0 for Claude subscriptions
- **Configuration**: JSON-based config system
- **Logging**: Structured logging with slog
- **Build**: Standard Go toolchain
- **Task Runner**: Taskfile.dev

## Architecture
```
crush/
├── internal/
│   ├── cmd/           # CLI commands (root, run, auth, logs, schema)
│   ├── app/           # Main application logic
│   ├── tui/           # Terminal user interface
│   ├── auth/          # OAuth authentication system
│   ├── llm/           # LLM provider implementations
│   ├── config/        # Configuration management
│   ├── db/            # Database operations
│   ├── session/       # Session management
│   └── ...
├── configs/           # Example configuration files
├── scripts/           # Build and utility scripts
└── main.go           # Application entry point
```

## Key New Features (Claude Subscription)
- OAuth authentication flow for Claude Pro/Max subscriptions
- `claudesub` provider for authenticated Claude access
- Token management and refresh capabilities
- Authentication status tracking
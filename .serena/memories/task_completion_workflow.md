# Task Completion Workflow

## Before Starting Development
1. **Environment Setup**
   ```bash
   # Ensure Go 1.25.0+ is installed
   go version
   
   # Install development dependencies
   task lint:install
   ```

2. **Build and Test**
   ```bash
   # Build the application
   go build -o ./crush .
   
   # Run existing tests to ensure baseline works
   task test
   ```

## During Development
1. **Code Quality Checks**
   ```bash
   # Run linter continuously during development
   task lint
   
   # Format code
   task fmt
   ```

2. **Testing New Features**
   ```bash
   # Test basic functionality
   ./crush --help
   ./crush auth status
   
   # Test with debug logging
   ./crush --debug run "test prompt"
   ```

## Task Completion Checklist
When a development task is completed, perform these steps in order:

### 1. Code Quality
```bash
# Format all code
task fmt

# Run linters and fix issues
task lint-fix

# Verify linting passes
task lint
```

### 2. Testing
```bash
# Run all tests
task test

# Test binary functionality
./crush --help
./crush auth status
./crush logs --tail 10

# Test new features specifically
# (add specific tests based on changes made)
```

### 3. Build Verification
```bash
# Clean build
go clean
go build -o ./crush .

# Verify binary works
./crush --version
```

### 4. Documentation
- Update README.md if user-facing changes
- Update configuration schema if config changes
- Add/update godoc comments for new/modified functions

### 5. Git Workflow
```bash
# Check status
git status

# Review changes
git diff

# Stage changes
git add .

# Commit with descriptive message
git commit -m "feat: add Claude subscription OAuth authentication"

# Push to branch
git push origin feature-branch
```

## Quality Gates
- [ ] All linters pass (`task lint`)
- [ ] All tests pass (`task test`)
- [ ] Code is properly formatted (`task fmt`)
- [ ] Binary builds successfully
- [ ] New features have appropriate tests
- [ ] Documentation is updated
- [ ] Git history is clean and descriptive

## Claude Subscription Specific Testing
When working on Claude subscription features:
```bash
# Test authentication flow
./crush auth status
./crush auth login  # (requires manual OAuth completion)
./crush auth logout

# Test provider functionality
CLAUDE_SUB_API_KEY="test" ./crush run "test prompt" --quiet

# Check logs for OAuth-related messages
./crush logs --tail 20 | grep -i oauth
```
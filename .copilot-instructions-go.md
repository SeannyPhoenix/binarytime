# Go Instructions for binarytime

## Project Structure
- **Main Package**: `pkg/binarytime/` - Core library implementation
- **CLI Tools**: `cmd/` - Command-line applications
- **Utilities**: `tools/` - Development and testing utilities
- **Internal**: `internal/` - Private packages not for external use

## Technology Stack
- **Language**: Go 1.21+ 
- **Dependencies**: Minimal external dependencies, prefer standard library
- **Testing**: Standard `go test` with table-driven tests
- **Build**: Standard `go build` and Makefile

## Code Standards
- Follow effective Go guidelines
- Use `gofmt` for consistent formatting
- Prefer composition over inheritance
- Use meaningful package and variable names
- Write comprehensive tests for all public APIs
- Document all exported functions and types

## Core Architecture
- **Fixed128**: 128-bit fixed-point arithmetic using `math/big.Int`
- **Date**: High-precision timestamp representation
- **Duration**: High-precision time span representation
- **Constants**: `dayNs = 86_400_000_000_000` (nanoseconds per day)

## Package Organization
```
pkg/
├── binarytime/         # Main library
│   ├── date.go         # Date functionality
│   ├── duration.go     # Duration functionality
│   ├── dateformat.go   # Formatting utilities
│   └── other.go        # Constants and utilities
├── fixed128/           # 128-bit arithmetic
├── byteglyph/          # Binary representation utilities
└── zordercurve/        # Z-order curve implementations
```

## Testing Guidelines
- Use table-driven tests for multiple test cases
- Test edge cases and error conditions
- Use `testing.T.Helper()` for test helper functions
- Group related tests in the same file
- Use descriptive test names that explain what is being tested

## Build and Development
- Use `make` for common build tasks
- Support cross-compilation for multiple platforms
- Maintain backward compatibility within major versions
- Use semantic versioning for releases

## Documentation Requirements
- Use Go doc comments for all exported APIs
- Include usage examples in package documentation
- Keep README.md updated with current functionality
- Document any breaking changes in release notes

## API Design Principles
- Prefer explicit error handling
- Use value receivers for immutable types
- Provide both convenience and precision methods
- Maintain consistency with Go standard library patterns
- Use interfaces for testability and flexibility
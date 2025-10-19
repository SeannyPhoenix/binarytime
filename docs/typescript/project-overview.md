# Project Overview

## Mission
Converting the Go-based `binarytime` library to a TypeScript NPM package for web and Node.js environments.

## Core Components
The original Go library provides high-precision time calculations using 128-bit fixed-point arithmetic:

- **Fixed128**: 128-bit fixed-point arithmetic using `big.Int`
- **Date**: Timestamp representation with nanosecond precision
- **Duration**: Time span representation
- **Constants**: `dayNs = 86_400_000_000_000` (nanoseconds per day)

## TypeScript Strategy

### Key Design Decisions
1. **Target Environment**: Browser + Node.js compatibility
2. **API Style**: JavaScript-idiomatic patterns (not direct Go translation)
3. **Precision**: Millisecond precision (JavaScript-appropriate)
4. **Testing Framework**: Jest with comprehensive coverage
5. **Build System**: ESBuild + TypeScript for fast builds

### Implementation Approach
- Replace Go's `big.Int` with JavaScript's native `BigInt`
- Maintain immutable API design from Go version
- Preserve high precision arithmetic operations
- Support both CommonJS and ES Module builds

## Project Phases
- ✅ **Phase 1**: NPM package setup and tooling configuration
- ✅ **Phase 2**: Core architecture (Fixed128, BinaryDate, BinaryDuration)
- ⏳ **Phase 3**: Advanced features (formatting, parsing, time zones)
- ⏳ **Phase 4**: Performance optimization and edge case handling
- ⏳ **Phase 5**: NPM publishing preparation
# NPM/TypeScript Instructions for @seannyphoenix/binarytime

## Package Configuration
- **Package Name**: `@seannyphoenix/binarytime`
- **Target Environments**: Both browser and Node.js (universal compatibility)
- **API Style**: JavaScript-idiomatic patterns (native feel for JS developers)
- **Source Directory**: `src/`
- **Build Directory**: `dist/`

## Technology Stack
- **Language**: TypeScript with strict configuration
- **Build Tool**: ESBuild for fast bundling
- **Testing**: Jest framework
- **Linting**: ESLint v9 with TypeScript rules
- **Package Type**: ES modules with CommonJS compatibility

## Code Standards
- Use TypeScript strict mode
- Follow ES2020+ standards
- Support BigInt for 128-bit arithmetic (replacing Go's big.Int)
- Maintain immutable API design from original Go version
- Use descriptive variable names and comprehensive JSDoc

## Build Outputs Required
- CommonJS bundle (`dist/index.js`)
- ES module bundle (`dist/index.esm.js`) 
- TypeScript declarations (`dist/index.d.ts`)

## Available NPM Scripts
- `npm run build`: Build all output formats
- `npm run dev`: Development watch mode
- `npm test`: Run Jest tests
- `npm run lint`: ESLint code quality checks
- `npm run lint:fix`: Auto-fix linting issues
- `npm run clean`: Remove dist folder

## Architecture Notes
- **Fixed128**: Core 128-bit fixed-point arithmetic using BigInt
- **BinaryDate**: High-precision timestamp representation
- **BinaryDuration**: High-precision time span representation
- **Constants**: `DAY_NANOSECONDS = 86_400_000_000_000n`

## File Organization
```
src/
├── __tests__/          # Jest test files
├── lib/                # Core library implementations
│   ├── fixed128.ts     # 128-bit arithmetic
│   ├── binary-date.ts  # Date functionality
│   └── binary-duration.ts # Duration functionality
└── index.ts            # Main export file
```

## Testing Guidelines
- Write comprehensive tests for all public APIs
- Test edge cases and error conditions
- Use Jest's built-in matchers when possible
- Group related tests in describe blocks

## Documentation Requirements
- Keep `docs/typescript-init.md` updated with process and decisions
- Use JSDoc for all public APIs
- Include usage examples in documentation
- Document any breaking changes from Go version

## Project Status
- ✅ Phase 1: NPM package setup and tooling configuration
- 🔄 Phase 2: Implement Fixed128 equivalent using BigInt
- ⏳ Phase 3: Port Date and Duration functionality
- ⏳ Phase 4: Comprehensive testing
- ⏳ Phase 5: NPM publishing preparation

## API Design Preferences
- Prefer JavaScript-idiomatic APIs over direct Go translations
- Use modern ES features where appropriate
- Maintain backward compatibility for widely-used features
- Prioritize developer experience and ease of use
- Support method chaining where it makes sense
- Use proper TypeScript generics for type safety
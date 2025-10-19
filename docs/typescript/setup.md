# Setup & Configuration

## NPM Package Configuration

**Package Name**: `@seannyphoenix/binarytime`
**Target Environments**: Browser + Node.js (universal compatibility)

### Dependencies Installed
```bash
npm install --save-dev typescript @types/node esbuild eslint @typescript-eslint/parser @typescript-eslint/eslint-plugin jest @types/jest ts-jest @eslint/js typescript-eslint
```

## Configuration Files

### package.json
- ES module type with CommonJS compatibility
- Build scripts for multiple output formats
- Proper scoped package naming
- Dev dependencies for TypeScript toolchain

### tsconfig.json
- Strict TypeScript configuration
- ES2020 target for modern JavaScript features
- Declaration file generation
- Source map support

### eslint.config.js
- ESLint v9 flat configuration
- TypeScript-specific rules
- Integration with @typescript-eslint
- Modern ES module format

### jest.config.cjs
- Jest testing framework configuration
- ts-jest integration for TypeScript support
- Test file patterns and coverage settings

## Build System

### ESBuild Configuration
- Fast bundling for development and production
- Multiple output formats:
  - CommonJS (`dist/index.js`)
  - ES Module (`dist/index.esm.js`)
  - TypeScript declarations (`dist/index.d.ts`)

### Available Scripts
- `npm run build`: Full production build
- `npm run dev`: Development watch mode
- `npm test`: Run Jest test suite
- `npm run lint`: ESLint code quality checks
- `npm run lint:fix`: Auto-fix linting issues
- `npm run clean`: Remove build output

## Directory Structure
```
src/
├── __tests__/          # Jest test files
│   ├── BinaryDate.test.ts
│   ├── BinaryDuration.test.ts
│   ├── Fixed128.test.ts
│   └── Fixed128-multiplication.test.ts
├── lib/                # Core library implementations
│   ├── Fixed128.ts     # 128-bit arithmetic
│   ├── BinaryDate.ts   # Date functionality
│   ├── BinaryDuration.ts # Duration functionality
│   └── constants.ts    # Shared constants
└── index.ts            # Main export file

dist/                   # Build output (generated)
├── index.js            # CommonJS bundle
├── index.esm.js        # ES module bundle
└── index.d.ts          # TypeScript declarations
```
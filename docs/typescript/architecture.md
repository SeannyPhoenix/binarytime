# Architecture

## Core Design Principles

### Immutability
All classes follow immutable patterns - operations return new instances rather than modifying existing ones.

### JavaScript Idioms
APIs are designed to feel natural to JavaScript developers rather than being direct Go translations.

### Precision Strategy
- **Target**: Millisecond precision (appropriate for JavaScript environments)
- **Storage**: BigInt for high-precision arithmetic
- **Compatibility**: JavaScript Date object interoperability

## Class Architecture

### Fixed128
**Purpose**: 128-bit fixed-point arithmetic using JavaScript BigInt

**Key Features**:
- Immutable value semantics
- Complete arithmetic operations (add, subtract, multiply, divide)
- Comparison operations
- Byte conversion utilities
- Component access (high/low 64-bit parts)

**TypeScript Interface**:
```typescript
class Fixed128 {
  constructor(high: bigint, low: bigint)
  static fromBigInt(value: bigint): Fixed128
  add(other: Fixed128): Fixed128
  subtract(other: Fixed128): Fixed128
  multiply(other: Fixed128): Fixed128
  divide(other: Fixed128): Fixed128
  // ... additional methods
}
```

### BinaryDate
**Purpose**: High-precision timestamp representation

**Key Features**:
- Creation from Unix timestamps (milliseconds, seconds)
- Conversion to/from JavaScript Date objects
- Arithmetic operations with time units
- Comparison operations
- Immutable design

**TypeScript Interface**:
```typescript
class BinaryDate {
  static fromUnixMillis(millis: number): BinaryDate
  static fromDate(date: Date): BinaryDate
  toDate(): Date
  addMilliseconds(ms: number): BinaryDate
  addDays(days: number): BinaryDate
  before(other: BinaryDate): boolean
  // ... additional methods
}
```

### BinaryDuration
**Purpose**: High-precision time span representation

**Key Features**:
- Creation from various time units
- Unit conversion methods
- Arithmetic operations
- Sign operations
- Scalar multiplication/division

**TypeScript Interface**:
```typescript
class BinaryDuration {
  static fromMilliseconds(ms: number): BinaryDuration
  static fromDays(days: number): BinaryDuration
  toMilliseconds(): number
  add(other: BinaryDuration): BinaryDuration
  multiply(scalar: number): BinaryDuration
  isNegative(): boolean
  // ... additional methods
}
```

## Constants
- `DAY_NANOSECONDS = 86_400_000_000_000n`: Nanoseconds per day (BigInt)
- Consistent with Go version but adapted for JavaScript BigInt

## Testing Strategy
- **Framework**: Jest with TypeScript support
- **Coverage**: Comprehensive test suites for all public APIs
- **Focus Areas**: 
  - Immutability verification
  - Precision accuracy
  - Edge cases and error conditions
  - Cross-format conversion accuracy

## Build Output
- **CommonJS**: For Node.js compatibility
- **ES Modules**: For modern bundlers and browsers
- **TypeScript Declarations**: For full type safety
- **Source Maps**: For debugging support
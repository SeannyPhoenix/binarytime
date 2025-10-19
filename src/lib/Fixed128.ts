/**
 * Fixed128 represents a 128-bit fixed-point fractional number.
 * The top 64 bits represent the whole part, and the bottom 64 bits represent
 * the fractional part. Uses JavaScript's native BigInt for the underlying storage.
 * 
 * This class behaves as a value type - all operations return new instances
 * and the original instances remain unchanged (immutable).
 */
export class Fixed128 {
  readonly value: bigint;

  /**
   * Static constants for common values
   */
  static readonly ZERO = new Fixed128(0n, 1n);
  static readonly ONE = new Fixed128(1n, 1n);

  /**
   * Create a new Fixed128 representing the fraction x/y
   * @param x - numerator
   * @param y - denominator
   * @throws Error if y is zero (division by zero)
   */
  constructor(x: bigint, y: bigint) {
    if (y === 0n) {
      throw new Error('Division by zero');
    }

    const negX = x < 0n;
    const negY = y < 0n;
    const neg = negX !== negY;

    const absX = x < 0n ? -x : x;
    const absY = y < 0n ? -y : y;

    const [hi, lo] = Fixed128.getComponents(absX, absY);
    this.value = Fixed128.assemble(neg, hi, lo);
  }

  /**
   * Create a Fixed128 from a raw BigInt value (internal use)
   */
  static fromBigInt(value: bigint): Fixed128 {
    const instance = Object.create(Fixed128.prototype);
    instance.value = value;
    return instance;
  }

  /**
   * Create a copy of this Fixed128 (returns this since it's immutable)
   */
  copy(): Fixed128 {
    return this;
  }

  /**
   * Get the underlying BigInt value
   */
  getValue(): bigint {
    return this.value;
  }

  /**
   * Get the sign of the number (-1, 0, or 1)
   */
  sign(): number {
    if (this.value < 0n) return -1;
    if (this.value > 0n) return 1;
    return 0;
  }

  /**
   * Compare with another Fixed128
   * @returns -1 if this < other, 0 if equal, 1 if this > other
   */
  cmp(other: Fixed128): number {
    if (this.value < other.value) return -1;
    if (this.value > other.value) return 1;
    return 0;
  }

  /**
   * Check if this number is negative
   */
  isNeg(): boolean {
    return this.value < 0n;
  }

  /**
   * Check if this number is zero
   */
  isZero(): boolean {
    return this.value === 0n;
  }

  /**
   * Get the high and low 64-bit components
   * @returns [hi, lo] as bigints
   */
  hiLo(): [bigint, bigint] {
    const [, hi, lo] = Fixed128.disassemble(this);
    return [hi, lo];
  }

  /**
   * Get the number as a byte array (16 bytes, big-endian)
   */
  bytes(): Uint8Array {
    const bytes = new Uint8Array(16);
    let value = this.value;

    // Handle negative numbers with two's complement
    const isNegative = value < 0n;
    if (isNegative) {
      // Convert to positive for byte extraction
      value = -value;
    }

    // Extract bytes from right to left (little-endian internally)
    for (let i = 15; i >= 0; i--) {
      bytes[i] = Number(value & 0xffn);
      value = value >> 8n;
    }

    // Apply two's complement for negative numbers
    if (isNegative) {
      // Invert all bits
      for (let i = 0; i < 16; i++) {
        bytes[i] = (~bytes[i]!) & 0xff;
      }

      // Add 1
      let carry = 1;
      for (let i = 15; i >= 0 && carry > 0; i--) {
        const sum = bytes[i]! + carry;
        bytes[i] = sum & 0xff;
        carry = sum >> 8;
      }
    }

    return bytes;
  }

  /**
   * Add another Fixed128 to this one
   */
  add(other: Fixed128): Fixed128 {
    return Fixed128.fromBigInt(this.value + other.value);
  }

  /**
   * Subtract another Fixed128 from this one
   */
  sub(other: Fixed128): Fixed128 {
    return Fixed128.fromBigInt(this.value - other.value);
  }

  /**
   * Multiply this Fixed128 by another
   */
  mul(other: Fixed128): Fixed128 {
    return Fixed128.fromBigInt(this.value * other.value);
  }

  /**
   * Divide this Fixed128 by another
   */
  quo(other: Fixed128): Fixed128 {
    if (other.value === 0n) {
      throw new Error('Division by zero');
    }
    return Fixed128.fromBigInt(this.value / other.value);
  }

  /**
   * Multiply by a bigint value and return the result as bigint
   */
  mulBigInt(y: bigint): bigint {
    // Short-circuit optimization for zero
    if (y === 0n) {
      return 0n;
    }

    const [negX, hi, lo] = Fixed128.disassemble(this);
    const negY = y < 0n;
    const absY = y < 0n ? -y : y;

    const whole = hi * absY;
    const part = Fixed128.hydrate(lo, absY);
    let result = whole + part;

    if (negX !== negY) {
      result = -result;
    }

    return result;
  }

  /**
   * Calculate components for x/y division
   */
  private static getComponents(x: bigint, y: bigint): [bigint, bigint] {
    if (y === 0n) {
      throw new Error('Division by zero in getComponents');
    }

    const hi = x / y;
    let part = x % y;

    const shift = Fixed128.leadingZeros64(y);
    y = y << BigInt(shift);
    part = part << BigInt(shift);

    let lo = 0n;
    let i = 0;

    while (i < 64 && y > 1n && part > 0n) {
      y = y >> 1n;
      const bit = part / y;
      part = part - (bit * y);
      lo = lo << 1n;
      lo = lo | bit;
      i++;
    }

    lo = lo << BigInt(64 - i);

    return [hi, lo];
  }

  /**
   * Assemble a Fixed128 from sign and components
   */
  private static assemble(neg: boolean, hi: bigint, lo: bigint): bigint {
    // Combine hi and lo into a 128-bit value
    let value = (hi << 64n) | lo;

    if (neg) {
      value = -value;
    }

    return value;
  }

  /**
   * Disassemble a Fixed128 into sign and components
   */
  private static disassemble(f128: Fixed128): [boolean, bigint, bigint] {
    const neg = f128.isNeg();
    const absValue = neg ? -f128.value : f128.value;

    const hi = absValue >> 64n;
    const lo = absValue & ((1n << 64n) - 1n);

    return [neg, hi, lo];
  }

  /**
   * Hydrate fractional part
   */
  private static hydrate(lo: bigint, div: bigint): bigint {
    const shift = Fixed128.leadingZeros64(div);
    div = div << BigInt(shift);

    let part = 0n;
    for (let i = 0; i < 64 && div > 0n; i++) {
      div = div >> 1n;
      const bit = (lo >> BigInt(63 - i)) & 1n;
      part = part + (div * bit);
    }

    part = Fixed128.round(shift, part);
    return part;
  }

  /**
   * Round the result based on shift
   */
  private static round(shift: number, part: bigint): bigint {
    if (shift === 0) {
      return part;
    }

    part = part >> BigInt(shift - 1);
    const bit = part & 1n;
    part = part >> 1n;
    part = part + bit;
    return part;
  }

  /**
   * Count leading zeros in a 64-bit value
   */
  private static leadingZeros64(x: bigint): number {
    if (x === 0n) return 64;

    let count = 0;
    let mask = 1n << 63n;

    for (let i = 0; i < 64; i++) {
      if ((x & mask) !== 0n) break;
      count++;
      mask = mask >> 1n;
    }

    return count;
  }

  /**
   * String representation for debugging
   */
  toString(): string {
    return `Fixed128(${this.value})`;
  }

  /**
   * Value equality check
   */
  equals(other: Fixed128): boolean {
    return this.value === other.value;
  }
}
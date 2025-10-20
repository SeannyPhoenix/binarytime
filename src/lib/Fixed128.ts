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
   * Get hex string representation with automatic precision
   * Similar to Go's String() method - shows minimal precision needed
   */
  toHexString(): string {
    if (this.isZero()) {
      return '00.00';
    }

    const bytes = this.bytesWithSign();

    // Find the first non-zero byte in the high part (skip sign byte)
    let high = 1;
    for (high = 1; high < 9 && bytes[high] === 0; high++) {
      // Continue until we find non-zero byte
    }

    // Find the last non-zero byte in the low part
    let low = 17;
    for (low = 17; low > 9 && bytes[low - 1] === 0; low--) {
      // Continue until we find non-zero byte
    }

    // Ensure we have at least minimal precision
    if (high >= 9) high = 8;
    if (low <= 9) low = 10;

    return this.toHexStringWithPrecision(high, low);
  }  /**
   * Get hex string representation with specified precision
   * @param high - starting byte index for high part (1-8)
   * @param low - ending byte index for low part (10-17)
   */
  toHexStringWithPrecision(high: number, low: number): string {
    const bytes = this.bytesWithSign();

    if (high >= 9 || low <= 9) {
      throw new Error(`Invalid precision: high=${high}, low=${low}`);
    }

    let result = '';

    // Add negative sign if needed
    if (bytes[0] === 1) {
      result += '-';
    }

    // High part (whole number)
    const highBytes = bytes.slice(high, 9);
    result += this.bytesToHex(highBytes);

    // Decimal point
    result += '.';

    // Low part (fractional)
    const lowBytes = bytes.slice(9, low);
    result += this.bytesToHex(lowBytes);

    return result;
  }

  /**
   * Get base64 string representation
   */
  toBase64(): string {
    const bytes = this.bytesWithSign();
    return btoa(String.fromCharCode(...bytes));
  }

  /**
   * Create Fixed128 from base64 string
   */
  static fromBase64(base64: string): Fixed128 {
    try {
      const binaryString = atob(base64);
      const bytes = new Uint8Array(binaryString.length);
      for (let i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
      }
      return Fixed128.fromBytesWithSign(bytes);
    } catch (error) {
      throw new Error(`Invalid base64 string: ${error}`);
    }
  }

  /**
   * Parse Fixed128 from hex string (e.g., "1A2B.3C4D" or "-1A2B.3C4D")
   */
  static fromHexString(hexStr: string): Fixed128 {
    if (!hexStr || hexStr.length === 0) {
      throw new Error('Empty hex string');
    }

    let str = hexStr.trim();
    const isNegative = str.startsWith('-');
    if (isNegative) {
      str = str.slice(1);
    }

    const parts = str.split('.');
    if (parts.length !== 2) {
      throw new Error('Expected format: "HI.LO"');
    }

    const [hiStr, loStr] = parts as [string, string];

    // Convert hex strings to bytes
    const hiBytes = Fixed128.hexToBytes(hiStr);
    const loBytes = Fixed128.hexToBytes(loStr);

    if (hiBytes.length > 8 || loBytes.length > 8) {
      throw new Error('Hex string too wide');
    }

    // Create 17-byte array (1 sign + 16 data bytes)
    const bytes = new Uint8Array(17);

    // Copy high bytes to positions 1-8 (right-aligned)
    bytes.set(hiBytes, 9 - hiBytes.length);

    // Copy low bytes to positions 9-16 (left-aligned)
    bytes.set(loBytes, 9);

    // Set sign
    if (isNegative) {
      bytes[0] = 1;
    }

    return Fixed128.fromBytesWithSign(bytes);
  }

  /**
   * Get bytes with sign byte (17 bytes total: 1 sign + 16 data)
   * Compatible with Go implementation
   */
  private bytesWithSign(): Uint8Array {
    const result = new Uint8Array(17);

    // For positive numbers, just use the value directly
    // For negative numbers, store the absolute value and set sign byte
    const isNegative = this.isNeg();
    const absValue = isNegative ? -this.value : this.value;

    // Set sign byte
    result[0] = isNegative ? 1 : 0;

    // Convert absolute value to bytes (big-endian)
    let value = absValue;
    for (let i = 16; i >= 1; i--) {
      result[i] = Number(value & 0xffn);
      value = value >> 8n;
    }

    return result;
  }

  /**
   * Create Fixed128 from bytes with sign (17 bytes: 1 sign + 16 data)
   */
  static fromBytesWithSign(bytes: Uint8Array): Fixed128 {
    if (bytes.length !== 17) {
      throw new Error(`Invalid length: expected 17, got ${bytes.length}`);
    }

    const signByte = bytes[0];
    if (signByte !== 0 && signByte !== 1) {
      throw new Error(`Invalid sign byte: ${signByte}`);
    }

    const dataBytes = bytes.slice(1);
    let value = 0n;

    // Convert bytes to BigInt (big-endian)
    for (let i = 0; i < dataBytes.length; i++) {
      value = (value << 8n) | BigInt(dataBytes[i]!);
    }

    // Apply sign
    if (signByte === 1) {
      value = -value;
    }

    return Fixed128.fromBigInt(value);
  }

  /**
   * Convert bytes to hex string
   */
  private bytesToHex(bytes: Uint8Array): string {
    return Array.from(bytes, byte => byte.toString(16).padStart(2, '0')).join('').toUpperCase();
  }

  /**
   * Convert hex string to bytes
   */
  private static hexToBytes(hex: string): Uint8Array {
    if (hex.length % 2 !== 0) {
      hex = '0' + hex; // Pad with leading zero
    }

    const bytes = new Uint8Array(hex.length / 2);
    for (let i = 0; i < hex.length; i += 2) {
      bytes[i / 2] = parseInt(hex.slice(i, i + 2), 16);
    }

    return bytes;
  }

  /**
   * Value equality check
   */
  equals(other: Fixed128): boolean {
    return this.value === other.value;
  }
}
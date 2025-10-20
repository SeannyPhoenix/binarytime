import { Fixed128 } from './Fixed128';
import { DAY_MILLISECONDS } from './constants';

/**
 * BinaryDuration represents a high-precision time span using 128-bit fixed-point arithmetic.
 * It provides millisecond precision and is immutable - all operations return new instances.
 */
export class BinaryDuration {
  private readonly value: Fixed128;

  /**
   * Create a new BinaryDuration from a Fixed128 value
   */
  constructor(value: Fixed128) {
    this.value = value;
  }

  /**
   * Create a BinaryDuration from milliseconds
   */
  static fromMillis(millis: bigint): BinaryDuration {
    const value = new Fixed128(millis, DAY_MILLISECONDS);
    return new BinaryDuration(value);
  }

  /**
   * Create a BinaryDuration from seconds
   */
  static fromSeconds(seconds: bigint): BinaryDuration {
    const millis = seconds * 1_000n;
    return BinaryDuration.fromMillis(millis);
  }

  /**
   * Create a BinaryDuration from minutes
   */
  static fromMinutes(minutes: bigint): BinaryDuration {
    const millis = minutes * 60_000n;
    return BinaryDuration.fromMillis(millis);
  }

  /**
   * Create a BinaryDuration from hours
   */
  static fromHours(hours: bigint): BinaryDuration {
    const millis = hours * 3_600_000n;
    return BinaryDuration.fromMillis(millis);
  }

  /**
   * Create a BinaryDuration from days
   */
  static fromDays(days: bigint): BinaryDuration {
    const millis = days * DAY_MILLISECONDS;
    return BinaryDuration.fromMillis(millis);
  }

  /**
   * Get duration in milliseconds
   */
  millis(): bigint {
    return this.value.mulBigInt(DAY_MILLISECONDS);
  }

  /**
   * Get duration in seconds (may lose precision)
   */
  seconds(): bigint {
    return this.millis() / 1_000n;
  }

  /**
   * Get duration in minutes (may lose precision)
   */
  minutes(): bigint {
    return this.millis() / 60_000n;
  }

  /**
   * Get duration in hours (may lose precision)
   */
  hours(): bigint {
    return this.millis() / 3_600_000n;
  }

  /**
   * Get duration in days (may lose precision)
   */
  days(): bigint {
    return this.millis() / DAY_MILLISECONDS;
  }

  /**
   * Create a copy of this BinaryDuration (returns this since it's immutable)
   */
  copy(): BinaryDuration {
    return this;
  }

  /**
   * Check if this BinaryDuration is zero
   */
  isZero(): boolean {
    return this.value.isZero();
  }

  /**
   * Check if this duration is negative
   */
  isNegative(): boolean {
    return this.value.isNeg();
  }

  /**
   * Check if two BinaryDuration instances are equal
   */
  equals(other: BinaryDuration): boolean {
    return this.value.equals(other.value);
  }

  /**
   * Compare with another BinaryDuration
   * @returns -1 if this < other, 0 if equal, 1 if this > other
   */
  cmp(other: BinaryDuration): number {
    return this.value.cmp(other.value);
  }

  /**
   * Add another duration to this one
   */
  add(other: BinaryDuration): BinaryDuration {
    const newValue = this.value.add(other.value);
    return new BinaryDuration(newValue);
  }

  /**
   * Subtract another duration from this one
   */
  sub(other: BinaryDuration): BinaryDuration {
    const newValue = this.value.sub(other.value);
    return new BinaryDuration(newValue);
  }

  /**
   * Multiply this duration by a scalar
   */
  mulScalar(scalar: bigint): BinaryDuration {
    const scaledMillis = this.millis() * scalar;
    return BinaryDuration.fromMillis(scaledMillis);
  }

  /**
   * Divide this duration by a scalar
   */
  divScalar(scalar: bigint): BinaryDuration {
    if (scalar === 0n) {
      throw new Error('Division by zero');
    }
    const dividedMillis = this.millis() / scalar;
    return BinaryDuration.fromMillis(dividedMillis);
  }

  /**
   * Get the absolute value of this duration
   */
  abs(): BinaryDuration {
    if (this.isNegative()) {
      return BinaryDuration.fromMillis(-this.millis());
    }
    return this;
  }

  /**
   * Negate this duration
   */
  neg(): BinaryDuration {
    return BinaryDuration.fromMillis(-this.millis());
  }

  /**
   * Get the underlying Fixed128 value
   */
  fixed128(): Fixed128 {
    return this.value;
  }

  /**
   * Get the underlying BigInt value
   */
  bigInt(): bigint {
    return this.value.getValue();
  }

  /**
   * Get the duration as a byte array (16 bytes, big-endian)
   */
  bytes(): Uint8Array {
    return this.value.bytes();
  }

  /**
   * String representation for debugging
   */
  toString(): string {
    const millis = this.millis();

    // Zero duration
    if (millis === 0n) return '0ms';

    const absMillis = millis < 0n ? -millis : millis;
    const sign = millis < 0n ? '-' : '';

    // Convert to human-readable format similar to Go's time.Duration
    const parts: string[] = [];

    // Days (not in Go's standard format, but useful for long durations)
    if (absMillis >= DAY_MILLISECONDS) {
      const days = absMillis / DAY_MILLISECONDS;
      const remainder = absMillis % DAY_MILLISECONDS;
      parts.push(`${days}d`);

      if (remainder > 0n) {
        return sign + parts.join('') + new BinaryDuration(new Fixed128(remainder, DAY_MILLISECONDS)).formatRemainder();
      }
      return sign + parts.join('');
    }

    return sign + this.formatRemainder();
  }

  /**
   * Format the remainder of a duration (less than a day)
   */
  private formatRemainder(): string {
    const absMillis = this.millis() < 0n ? -this.millis() : this.millis();
    const parts: string[] = [];

    let remaining = absMillis;

    // Hours
    if (remaining >= 3_600_000n) {
      const hours = remaining / 3_600_000n;
      remaining = remaining % 3_600_000n;
      parts.push(`${hours}h`);
    }

    // Minutes
    if (remaining >= 60_000n) {
      const minutes = remaining / 60_000n;
      remaining = remaining % 60_000n;
      parts.push(`${minutes}m`);
    }

    // Seconds and milliseconds
    if (remaining >= 1_000n) {
      const seconds = remaining / 1_000n;
      const ms = remaining % 1_000n;
      if (ms === 0n) {
        parts.push(`${seconds}s`);
      } else {
        // Format as decimal seconds
        const secWithMs = Number(seconds) + Number(ms) / 1000;
        parts.push(`${secWithMs}s`);
      }
    } else if (remaining > 0n) {
      // Just milliseconds
      parts.push(`${remaining}ms`);
    }

    return parts.join('') || '0ms';
  }

  /**
   * Parse a duration string in Go time.Duration format
   * Supports: "300ms", "1.5h", "2h45m", "1d2h30m", etc.
   * Valid units: "ms", "s", "m", "h", "d"
   */
  static parse(durationStr: string): BinaryDuration {
    if (!durationStr || durationStr.trim().length === 0) {
      throw new Error('Empty duration string');
    }

    let str = durationStr.trim();
    const isNegative = str.startsWith('-');
    if (isNegative) {
      str = str.slice(1);
    }

    let totalMillis = 0n;

    // Regex to match number+unit pairs
    const unitRegex = /([0-9]*\.?[0-9]+)([a-zA-Z]+)/g;
    let match;
    let hasMatches = false;

    while ((match = unitRegex.exec(str)) !== null) {
      hasMatches = true;
      const valueStr = match[1]!;
      const unit = match[2]!.toLowerCase();

      // Parse the numeric value (may be decimal)
      const value = parseFloat(valueStr);
      if (isNaN(value)) {
        throw new Error(`Invalid number: ${valueStr}`);
      }

      // Convert to milliseconds based on unit
      let multiplier: number;
      switch (unit) {
        case 'ms':
          multiplier = 1;
          break;
        case 's':
          multiplier = 1_000;
          break;
        case 'm':
          multiplier = 60_000;
          break;
        case 'h':
          multiplier = 3_600_000;
          break;
        case 'd':
          multiplier = Number(DAY_MILLISECONDS);
          break;
        default:
          throw new Error(`Unknown time unit: ${unit}`);
      }

      const millis = Math.round(value * multiplier);
      totalMillis += BigInt(millis);
    }

    if (!hasMatches) {
      throw new Error(`Invalid duration format: ${durationStr}`);
    }

    // Check if we parsed the entire string
    const reconstructed = str.replace(/([0-9]*\.?[0-9]+)([a-zA-Z]+)/g, '');
    if (reconstructed.trim().length > 0) {
      throw new Error(`Invalid duration format: ${durationStr}`);
    }

    const result = BinaryDuration.fromMillis(totalMillis);
    return isNegative ? result.neg() : result;
  }

  /**
   * Get hex string representation
   */
  hex(): string {
    return this.value.toHexString();
  }

  /**
   * Get hex string representation with fine precision
   */
  hexFine(): string {
    return this.value.toHexString();
  }

  /**
   * Get base64 string representation
   */
  base64(): string {
    return this.value.toBase64();
  }

  /**
   * Create BinaryDuration from hex string representation
   */
  static fromHex(hexStr: string): BinaryDuration {
    const value = Fixed128.fromHexString(hexStr);
    return new BinaryDuration(value);
  }

  /**
   * Create BinaryDuration from base64 string representation
   */
  static fromBase64(base64Str: string): BinaryDuration {
    const value = Fixed128.fromBase64(base64Str);
    return new BinaryDuration(value);
  }
}
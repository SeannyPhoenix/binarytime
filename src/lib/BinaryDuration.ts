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

    // Convert to a human-readable format
    if (millis === 0n) return 'BinaryDuration(0ms)';

    const absMillis = millis < 0n ? -millis : millis;
    const sign = millis < 0n ? '-' : '';

    // Choose appropriate unit
    if (absMillis >= DAY_MILLISECONDS) {
      const days = absMillis / DAY_MILLISECONDS;
      return `BinaryDuration(${sign}${days}d)`;
    } else if (absMillis >= 3600000n) { // 1 hour in ms
      const hours = absMillis / 3600000n;
      return `BinaryDuration(${sign}${hours}h)`;
    } else if (absMillis >= 60000n) { // 1 minute in ms
      const minutes = absMillis / 60000n;
      return `BinaryDuration(${sign}${minutes}m)`;
    } else if (absMillis >= 1000n) { // 1 second in ms
      const seconds = absMillis / 1000n;
      return `BinaryDuration(${sign}${seconds}s)`;
    } else {
      return `BinaryDuration(${sign}${absMillis}ms)`;
    }
  }
}
import { Fixed128 } from './Fixed128';
import { DAY_MILLISECONDS } from './constants';

/**
 * BinaryDate represents a high-precision timestamp using 128-bit fixed-point arithmetic.
 * It provides millisecond precision and is immutable - all operations return new instances.
 */
export class BinaryDate {
  private readonly value: Fixed128;

  /**
   * Create a new BinaryDate from a Fixed128 value
   */
  constructor(value: Fixed128) {
    this.value = value;
  }

  /**
   * Create a BinaryDate representing the current time
   */
  static now(): BinaryDate {
    return BinaryDate.fromDate(new Date());
  }

  /**
   * Create a BinaryDate from a JavaScript Date object
   */
  static fromDate(date: Date): BinaryDate {
    const millis = BigInt(date.getTime());
    return BinaryDate.fromUnixMillis(millis);
  }

  /**
   * Create a BinaryDate from Unix timestamp in milliseconds
   */
  static fromUnixMillis(millis: bigint): BinaryDate {
    const value = new Fixed128(millis, DAY_MILLISECONDS);
    return new BinaryDate(value);
  }

  /**
   * Create a BinaryDate from Unix timestamp in seconds
   */
  static fromUnixSeconds(seconds: bigint): BinaryDate {
    const millis = seconds * 1_000n;
    return BinaryDate.fromUnixMillis(millis);
  }

  /**
   * Convert to JavaScript Date object
   */
  toDate(): Date {
    const millis = this.unixMilli();
    return new Date(Number(millis));
  }

  /**
   * Get Unix timestamp in milliseconds
   */
  unixMilli(): bigint {
    return this.value.mulBigInt(DAY_MILLISECONDS);
  }

  /**
   * Get Unix timestamp in seconds
   */
  unixSecond(): bigint {
    return this.unixMilli() / 1_000n;
  }

  /**
   * Create a copy of this BinaryDate (returns this since it's immutable)
   */
  copy(): BinaryDate {
    return this;
  }

  /**
   * Check if this BinaryDate is zero (Unix epoch)
   */
  isZero(): boolean {
    return this.value.isZero();
  }

  /**
   * Check if two BinaryDate instances are equal
   */
  equals(other: BinaryDate): boolean {
    return this.value.equals(other.value);
  }

  /**
   * Compare with another BinaryDate
   * @returns -1 if this < other, 0 if equal, 1 if this > other
   */
  cmp(other: BinaryDate): number {
    return this.value.cmp(other.value);
  }

  /**
   * Check if this date is before another date
   */
  before(other: BinaryDate): boolean {
    return this.cmp(other) < 0;
  }

  /**
   * Check if this date is after another date
   */
  after(other: BinaryDate): boolean {
    return this.cmp(other) > 0;
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
   * Get the date as a byte array (16 bytes, big-endian)
   */
  bytes(): Uint8Array {
    return this.value.bytes();
  }

  /**
   * Add milliseconds to this date
   */
  addMillis(millis: bigint): BinaryDate {
    const duration = new Fixed128(millis, DAY_MILLISECONDS);
    const newValue = this.value.add(duration);
    return new BinaryDate(newValue);
  }

  /**
   * Subtract milliseconds from this date
   */
  subMillis(millis: bigint): BinaryDate {
    const duration = new Fixed128(millis, DAY_MILLISECONDS);
    const newValue = this.value.sub(duration);
    return new BinaryDate(newValue);
  }

  /**
   * Add seconds to this date
   */
  addSeconds(seconds: bigint): BinaryDate {
    return this.addMillis(seconds * 1_000n);
  }

  /**
   * Subtract seconds from this date
   */
  subSeconds(seconds: bigint): BinaryDate {
    return this.subMillis(seconds * 1_000n);
  }

  /**
   * Add minutes to this date
   */
  addMinutes(minutes: bigint): BinaryDate {
    return this.addMillis(minutes * 60_000n);
  }

  /**
   * Subtract minutes from this date
   */
  subMinutes(minutes: bigint): BinaryDate {
    return this.subMillis(minutes * 60_000n);
  }

  /**
   * Add hours to this date
   */
  addHours(hours: bigint): BinaryDate {
    return this.addMillis(hours * 3_600_000n);
  }

  /**
   * Subtract hours from this date
   */
  subHours(hours: bigint): BinaryDate {
    return this.subMillis(hours * 3_600_000n);
  }

  /**
   * Add days to this date
   */
  addDays(days: bigint): BinaryDate {
    return this.addMillis(days * DAY_MILLISECONDS);
  }

  /**
   * Subtract days from this date
   */
  subDays(days: bigint): BinaryDate {
    return this.subMillis(days * DAY_MILLISECONDS);
  }

  /**
   * String representation for debugging
   */
  toString(): string {
    const date = this.toDate();
    return `BinaryDate(${date.toISOString()})`;
  }

  /**
   * Get ISO string representation
   */
  toISOString(): string {
    return this.toDate().toISOString();
  }
}
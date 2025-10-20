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
   * Add a duration (Fixed128) to this date
   * @param duration - the duration as Fixed128 to add
   */
  addFixed128Duration(duration: Fixed128): BinaryDate {
    const newValue = this.value.add(duration);
    return new BinaryDate(newValue);
  }

  /**
   * Subtract a duration (Fixed128) from this date
   * @param duration - the duration as Fixed128 to subtract
   */
  subFixed128Duration(duration: Fixed128): BinaryDate {
    const newValue = this.value.sub(duration);
    return new BinaryDate(newValue);
  }

  /**
   * Calculate the duration between this date and another date as Fixed128
   * @param other - the other date
   * @returns positive duration if other is after this date, negative if before
   */
  durationUntilFixed128(other: BinaryDate): Fixed128 {
    return other.value.sub(this.value);
  }

  /**
   * Calculate the duration since another date as Fixed128
   * @param other - the other date  
   * @returns positive duration if this date is after other, negative if before
   */
  durationSinceFixed128(other: BinaryDate): Fixed128 {
    return this.value.sub(other.value);
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

  /**
   * Get hex string representation with automatic precision
   * Equivalent to Go's Hex() method
   */
  hex(): string {
    return this.hexWithPrecision(8, 10);
  }

  /**
   * Get hex string representation with fine precision
   * Equivalent to Go's HexFine() method  
   */
  hexFine(): string {
    return this.value.toHexString();
  }

  /**
   * Get hex string representation with specified precision
   * @param high - starting byte index for high part (1-8)
   * @param low - ending byte index for low part (10-17)
   */
  hexWithPrecision(high: number, low: number): string {
    return this.value.toHexStringWithPrecision(high, low);
  }

  /**
   * Get base64 string representation
   * Equivalent to Go's Base64() method
   */
  base64(): string {
    return this.value.toBase64();
  }

  /**
   * Create BinaryDate from hex string representation
   * @param hexStr - hex string in format "HI.LO" or "-HI.LO"
   */
  static fromHex(hexStr: string): BinaryDate {
    const value = Fixed128.fromHexString(hexStr);
    return new BinaryDate(value);
  }

  /**
   * Create BinaryDate from base64 string representation
   */
  static fromBase64(base64Str: string): BinaryDate {
    const value = Fixed128.fromBase64(base64Str);
    return new BinaryDate(value);
  }
}
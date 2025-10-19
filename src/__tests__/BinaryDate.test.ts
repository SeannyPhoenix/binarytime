import { BinaryDate } from '../lib/BinaryDate';
import { Fixed128 } from '../lib/Fixed128';

describe('BinaryDate', () => {
  describe('creation methods', () => {
    it('should create from current time', () => {
      const date = BinaryDate.now();
      expect(date).toBeInstanceOf(BinaryDate);
      expect(date.isZero()).toBe(false);
    });

    it('should create from JavaScript Date', () => {
      const jsDate = new Date('2023-01-01T00:00:00.000Z');
      const binaryDate = BinaryDate.fromDate(jsDate);

      expect(binaryDate).toBeInstanceOf(BinaryDate);
      expect(binaryDate.toDate().getTime()).toBe(jsDate.getTime());
    });

    it('should create from Unix milliseconds', () => {
      const millis = 1672531200000n; // 2023-01-01T00:00:00.000Z in milliseconds
      const date = BinaryDate.fromUnixMillis(millis);

      expect(date.unixMilli()).toBe(millis);
    }); it('should create from Unix seconds', () => {
      const seconds = 1672531200n; // 2023-01-01T00:00:00.000Z in seconds
      const date = BinaryDate.fromUnixSeconds(seconds);

      expect(date.unixSecond()).toBe(seconds);
    });
  });

  describe('conversion methods', () => {
    it('should convert to JavaScript Date', () => {
      const original = new Date('2023-01-01T12:30:45.123Z');
      const binaryDate = BinaryDate.fromDate(original);
      const converted = binaryDate.toDate();

      expect(converted.getTime()).toBe(original.getTime());
    });

    it('should convert to Unix timestamps', () => {
      const millis = 1672574445123n;
      const date = BinaryDate.fromUnixMillis(millis);

      expect(date.unixMilli()).toBe(millis);
      expect(date.unixSecond()).toBe(millis / 1_000n);
    });

    it('should convert to ISO string', () => {
      const jsDate = new Date('2023-01-01T12:30:45.123Z');
      const binaryDate = BinaryDate.fromDate(jsDate);

      expect(binaryDate.toISOString()).toBe(jsDate.toISOString());
    });
  }); describe('comparison operations', () => {
    it('should compare dates correctly', () => {
      const earlier = BinaryDate.fromUnixSeconds(1672531200n); // 2023-01-01T00:00:00Z
      const later = BinaryDate.fromUnixSeconds(1672617600n);   // 2023-01-02T00:00:00Z

      expect(earlier.cmp(later)).toBe(-1);
      expect(later.cmp(earlier)).toBe(1);
      expect(earlier.cmp(earlier)).toBe(0);
    });

    it('should check before/after correctly', () => {
      const earlier = BinaryDate.fromUnixSeconds(1672531200n);
      const later = BinaryDate.fromUnixSeconds(1672617600n);

      expect(earlier.before(later)).toBe(true);
      expect(later.before(earlier)).toBe(false);
      expect(earlier.after(later)).toBe(false);
      expect(later.after(earlier)).toBe(true);
    });

    it('should check equality correctly', () => {
      const date1 = BinaryDate.fromUnixSeconds(1672531200n);
      const date2 = BinaryDate.fromUnixSeconds(1672531200n);
      const date3 = BinaryDate.fromUnixSeconds(1672617600n);

      expect(date1.equals(date2)).toBe(true);
      expect(date1.equals(date3)).toBe(false);
    });
  });

  describe('arithmetic operations', () => {
    it('should add milliseconds', () => {
      const date = BinaryDate.fromUnixSeconds(1672531200n);
      const added = date.addMillis(1000n); // Add 1 second

      expect(added.unixSecond()).toBe(1672531201n);
    });

    it('should subtract milliseconds', () => {
      const date = BinaryDate.fromUnixSeconds(1672531200n);
      const subtracted = date.subMillis(1000n); // Subtract 1 second

      expect(subtracted.unixSecond()).toBe(1672531199n);
    });

    it('should add seconds', () => {
      const date = BinaryDate.fromUnixSeconds(1672531200n);
      const added = date.addSeconds(60n); // Add 1 minute

      expect(added.unixSecond()).toBe(1672531260n);
    });

    it('should add minutes', () => {
      const date = BinaryDate.fromUnixSeconds(1672531200n);
      const added = date.addMinutes(60n); // Add 1 hour

      expect(added.unixSecond()).toBe(1672534800n);
    });

    it('should add hours', () => {
      const date = BinaryDate.fromUnixSeconds(1672531200n);
      const added = date.addHours(24n); // Add 1 day

      expect(added.unixSecond()).toBe(1672617600n);
    });

    it('should add days', () => {
      const date = BinaryDate.fromUnixSeconds(1672531200n);
      const added = date.addDays(1n); // Add 1 day

      expect(added.unixSecond()).toBe(1672617600n);
    });
  }); describe('immutability', () => {
    it('should not modify original date during operations', () => {
      const original = BinaryDate.fromUnixSeconds(1672531200n);
      const originalTime = original.unixSecond();

      original.addSeconds(60n);

      expect(original.unixSecond()).toBe(originalTime);
    });

    it('should return same instance from copy', () => {
      const original = BinaryDate.fromUnixSeconds(1672531200n);
      const copy = original.copy();

      expect(copy).toBe(original); // Same reference since immutable
    });
  });

  describe('zero value', () => {
    it('should detect zero date (Unix epoch)', () => {
      const epoch = BinaryDate.fromUnixSeconds(0n);
      expect(epoch.isZero()).toBe(true);

      const nonEpoch = BinaryDate.fromUnixSeconds(1n);
      expect(nonEpoch.isZero()).toBe(false);
    });
  });

  describe('data access', () => {
    it('should provide access to underlying Fixed128', () => {
      const date = BinaryDate.fromUnixSeconds(1672531200n);
      const fixed128 = date.fixed128();

      expect(fixed128).toBeInstanceOf(Fixed128);
    });

    it('should provide access to underlying BigInt', () => {
      const date = BinaryDate.fromUnixSeconds(1672531200n);
      const bigint = date.bigInt();

      expect(typeof bigint).toBe('bigint');
    });

    it('should convert to bytes', () => {
      const date = BinaryDate.fromUnixSeconds(1672531200n);
      const bytes = date.bytes();

      expect(bytes).toBeInstanceOf(Uint8Array);
      expect(bytes.length).toBe(16);
    });
  });

  describe('string representation', () => {
    it('should have toString method', () => {
      const date = BinaryDate.fromDate(new Date('2023-01-01T12:30:45.123Z'));
      const str = date.toString();

      expect(typeof str).toBe('string');
      expect(str).toContain('BinaryDate');
      expect(str).toContain('2023-01-01T12:30:45.123Z');
    });
  });
});
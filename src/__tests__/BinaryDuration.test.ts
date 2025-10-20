import { BinaryDuration } from '../lib/BinaryDuration';
import { Fixed128 } from '../lib/Fixed128';

describe('BinaryDuration', () => {
  describe('creation methods', () => {
    it('should create from milliseconds', () => {
      const duration = BinaryDuration.fromMillis(1000n); // 1 second
      expect(duration).toBeInstanceOf(BinaryDuration);
      expect(duration.millis()).toBe(1000n);
    });

    it('should create from seconds', () => {
      const duration = BinaryDuration.fromSeconds(60n); // 1 minute
      expect(duration.seconds()).toBe(60n);
      expect(duration.minutes()).toBe(1n);
    });

    it('should create from minutes', () => {
      const duration = BinaryDuration.fromMinutes(60n); // 1 hour
      expect(duration.minutes()).toBe(60n);
      expect(duration.hours()).toBe(1n);
    });

    it('should create from hours', () => {
      const duration = BinaryDuration.fromHours(24n); // 1 day
      expect(duration.hours()).toBe(24n);
      expect(duration.days()).toBe(1n);
    });

    it('should create from days', () => {
      const duration = BinaryDuration.fromDays(7n); // 1 week
      expect(duration.days()).toBe(7n);
    });
  });

  describe('conversion methods', () => {
    it('should convert between time units correctly', () => {
      const duration = BinaryDuration.fromSeconds(3661n); // 1 hour, 1 minute, 1 second

      expect(duration.millis()).toBe(3661000n);
      expect(duration.seconds()).toBe(3661n);
      expect(duration.minutes()).toBe(61n); // 61 minutes total
      expect(duration.hours()).toBe(1n); // 1 hour total (truncated)
    });

    it('should handle fractional conversions', () => {
      const duration = BinaryDuration.fromMillis(1500n); // 1.5 seconds

      expect(duration.millis()).toBe(1500n);
      expect(duration.seconds()).toBe(1n); // Truncated to 1 second
    });
  }); describe('arithmetic operations', () => {
    it('should add durations', () => {
      const duration1 = BinaryDuration.fromSeconds(30n);
      const duration2 = BinaryDuration.fromSeconds(45n);
      const sum = duration1.add(duration2);

      expect(sum.seconds()).toBe(75n);
    });

    it('should subtract durations', () => {
      const duration1 = BinaryDuration.fromMinutes(5n);
      const duration2 = BinaryDuration.fromMinutes(2n);
      const diff = duration1.sub(duration2);

      expect(diff.minutes()).toBe(3n);
    });

    it('should multiply by scalar', () => {
      const duration = BinaryDuration.fromSeconds(10n);
      const multiplied = duration.mulScalar(3n);

      expect(multiplied.seconds()).toBe(30n);
    });

    it('should divide by scalar', () => {
      const duration = BinaryDuration.fromSeconds(30n);
      const divided = duration.divScalar(3n);

      expect(divided.seconds()).toBe(10n);
    });

    it('should throw on division by zero', () => {
      const duration = BinaryDuration.fromSeconds(10n);
      expect(() => duration.divScalar(0n)).toThrow('Division by zero');
    });
  });

  describe('sign operations', () => {
    it('should handle negative durations', () => {
      const duration = BinaryDuration.fromSeconds(-30n);

      expect(duration.isNegative()).toBe(true);
      expect(duration.isZero()).toBe(false);
    });

    it('should calculate absolute value', () => {
      const negative = BinaryDuration.fromSeconds(-30n);
      const positive = negative.abs();

      expect(positive.isNegative()).toBe(false);
      expect(positive.seconds()).toBe(30n);
    });

    it('should negate duration', () => {
      const positive = BinaryDuration.fromSeconds(30n);
      const negative = positive.neg();

      expect(negative.isNegative()).toBe(true);
      expect(negative.seconds()).toBe(-30n);
    });

    it('should detect zero duration', () => {
      const zero = BinaryDuration.fromMillis(0n);
      const nonZero = BinaryDuration.fromMillis(1n);

      expect(zero.isZero()).toBe(true);
      expect(nonZero.isZero()).toBe(false);
    });
  });

  describe('comparison operations', () => {
    it('should compare durations correctly', () => {
      const shorter = BinaryDuration.fromSeconds(30n);
      const longer = BinaryDuration.fromMinutes(1n);

      expect(shorter.cmp(longer)).toBe(-1);
      expect(longer.cmp(shorter)).toBe(1);
      expect(shorter.cmp(shorter)).toBe(0);
    });

    it('should check equality correctly', () => {
      const duration1 = BinaryDuration.fromSeconds(60n);
      const duration2 = BinaryDuration.fromMinutes(1n);
      const duration3 = BinaryDuration.fromSeconds(30n);

      expect(duration1.equals(duration2)).toBe(true);
      expect(duration1.equals(duration3)).toBe(false);
    });
  });

  describe('immutability', () => {
    it('should not modify original duration during operations', () => {
      const original = BinaryDuration.fromSeconds(30n);
      const originalMillis = original.millis();

      original.add(BinaryDuration.fromSeconds(10n));

      expect(original.millis()).toBe(originalMillis);
    }); it('should return same instance from copy', () => {
      const original = BinaryDuration.fromSeconds(30n);
      const copy = original.copy();

      expect(copy).toBe(original); // Same reference since immutable
    });
  });

  describe('data access', () => {
    it('should provide access to underlying Fixed128', () => {
      const duration = BinaryDuration.fromSeconds(30n);
      const fixed128 = duration.fixed128();

      expect(fixed128).toBeInstanceOf(Fixed128);
    });

    it('should provide access to underlying BigInt', () => {
      const duration = BinaryDuration.fromSeconds(30n);
      const bigint = duration.bigInt();

      expect(typeof bigint).toBe('bigint');
    });

    it('should convert to bytes', () => {
      const duration = BinaryDuration.fromSeconds(30n);
      const bytes = duration.bytes();

      expect(bytes).toBeInstanceOf(Uint8Array);
      expect(bytes.length).toBe(16);
    });
  });

  describe('string representation', () => {
    it('should format zero duration', () => {
      const duration = BinaryDuration.fromMillis(0n);
      expect(duration.toString()).toBe('0ms');
    });

    it('should format milliseconds', () => {
      const duration = BinaryDuration.fromMillis(123n);
      expect(duration.toString()).toBe('123ms');
    });

    it('should format seconds', () => {
      const duration = BinaryDuration.fromSeconds(45n);
      expect(duration.toString()).toBe('45s');
    });

    it('should format minutes', () => {
      const duration = BinaryDuration.fromMinutes(30n);
      expect(duration.toString()).toBe('30m');
    });

    it('should format hours', () => {
      const duration = BinaryDuration.fromHours(12n);
      expect(duration.toString()).toBe('12h');
    });

    it('should format days', () => {
      const duration = BinaryDuration.fromDays(5n);
      expect(duration.toString()).toBe('5d');
    });

    it('should format negative durations', () => {
      const duration = BinaryDuration.fromSeconds(-30n);
      expect(duration.toString()).toBe('-30s');
    });
  });
});
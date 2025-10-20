import { BinaryDate, BinaryDuration, addDuration, subDuration, durationBetween, durationSince } from '../index';

describe('Duration Arithmetic with BinaryDate', () => {
  describe('addDuration function', () => {
    test('should add duration to date correctly', () => {
      const date = BinaryDate.fromUnixMillis(1640995200000n); // 2022-01-01
      const duration = BinaryDuration.fromHours(2n);

      const result = addDuration(date, duration);
      const expected = date.addMillis(2n * 3600000n);

      expect(result.equals(expected)).toBe(true);
    });

    test('should add negative duration (subtract)', () => {
      const date = BinaryDate.fromUnixMillis(1640995200000n);
      const duration = BinaryDuration.fromHours(-1n);

      const result = addDuration(date, duration);
      const expected = date.subMillis(3600000n);

      expect(result.equals(expected)).toBe(true);
    });
  });

  describe('subDuration function', () => {
    test('should subtract duration from date correctly', () => {
      const date = BinaryDate.fromUnixMillis(1640995200000n);
      const duration = BinaryDuration.fromMinutes(30n);

      const result = subDuration(date, duration);
      const expected = date.subMillis(30n * 60000n);

      expect(result.equals(expected)).toBe(true);
    });

    test('should subtract negative duration (add)', () => {
      const date = BinaryDate.fromUnixMillis(1640995200000n);
      const duration = BinaryDuration.fromMinutes(-30n);

      const result = subDuration(date, duration);
      const expected = date.addMillis(30n * 60000n);

      expect(result.equals(expected)).toBe(true);
    });
  });

  describe('durationBetween function', () => {
    test('should calculate positive duration between dates', () => {
      const date1 = BinaryDate.fromUnixMillis(1640995200000n); // 2022-01-01
      const date2 = BinaryDate.fromUnixMillis(1640995200000n + 3600000n); // 1 hour later

      const duration = durationBetween(date1, date2);

      expect(duration.hours()).toBe(1n);
      expect(duration.isNegative()).toBe(false);
    });

    test('should calculate negative duration between dates', () => {
      const date1 = BinaryDate.fromUnixMillis(1640995200000n + 3600000n); // 1 hour later
      const date2 = BinaryDate.fromUnixMillis(1640995200000n); // 2022-01-01

      const duration = durationBetween(date1, date2);

      expect(duration.hours()).toBe(-1n);
      expect(duration.isNegative()).toBe(true);
    });

    test('should return zero duration for same dates', () => {
      const date = BinaryDate.fromUnixMillis(1640995200000n);

      const duration = durationBetween(date, date);

      expect(duration.isZero()).toBe(true);
    });
  });

  describe('durationSince function', () => {
    test('should calculate duration since earlier date', () => {
      const earlier = BinaryDate.fromUnixMillis(1640995200000n); // 2022-01-01
      const later = BinaryDate.fromUnixMillis(1640995200000n + 7200000n); // 2 hours later

      const duration = durationSince(earlier, later);

      expect(duration.hours()).toBe(2n);
      expect(duration.isNegative()).toBe(false);
    });

    test('should calculate negative duration since later date', () => {
      const later = BinaryDate.fromUnixMillis(1640995200000n + 7200000n); // 2 hours later
      const earlier = BinaryDate.fromUnixMillis(1640995200000n); // 2022-01-01

      const duration = durationSince(later, earlier);

      expect(duration.hours()).toBe(-2n);
      expect(duration.isNegative()).toBe(true);
    });
  });

  describe('Complex duration arithmetic', () => {
    test('should handle chained operations', () => {
      const start = BinaryDate.fromUnixMillis(1640995200000n);
      const duration1 = BinaryDuration.fromHours(2n);
      const duration2 = BinaryDuration.fromMinutes(30n);

      // Add 2 hours, then subtract 30 minutes
      const intermediate = addDuration(start, duration1);
      const result = subDuration(intermediate, duration2);

      // Should be 1.5 hours after start
      const expected = start.addMillis(90n * 60000n);
      expect(result.equals(expected)).toBe(true);
    });

    test('should work with parsed duration strings', () => {
      const date = BinaryDate.fromUnixMillis(1640995200000n);
      const duration = BinaryDuration.parse('2h30m45s');

      const result = addDuration(date, duration);
      const expectedMillis = 2n * 3600000n + 30n * 60000n + 45n * 1000n;
      const expected = date.addMillis(expectedMillis);

      expect(result.equals(expected)).toBe(true);
    });

    test('should maintain precision through operations', () => {
      const date = BinaryDate.fromUnixMillis(1640995200123n); // with milliseconds
      const duration = BinaryDuration.fromMillis(456n);

      const result = addDuration(date, duration);
      const expectedMillis = 1640995200123n + 456n;

      expect(result.unixMilli()).toBe(expectedMillis);
    });
  });

  describe('Round-trip operations', () => {
    test('should support adding and subtracting same duration', () => {
      const original = BinaryDate.fromUnixMillis(1640995200000n);
      const duration = BinaryDuration.fromHours(5n);

      const added = addDuration(original, duration);
      const restored = subDuration(added, duration);

      expect(restored.equals(original)).toBe(true);
    });

    test('should calculate exact duration between dates', () => {
      const date1 = BinaryDate.fromUnixMillis(1640995200000n);
      const date2 = BinaryDate.fromUnixMillis(1640995200000n + 12345678n);

      const duration = durationBetween(date1, date2);
      const reconstructed = addDuration(date1, duration);

      expect(reconstructed.equals(date2)).toBe(true);
    });
  });
});
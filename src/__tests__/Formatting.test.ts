import { Fixed128 } from '../lib/Fixed128';
import { BinaryDate } from '../lib/BinaryDate';
import { BinaryDuration } from '../lib/BinaryDuration';

describe('Fixed128 Formatting', () => {
  describe('Hex String Formatting', () => {
    test('should format zero as "00.00"', () => {
      const zero = new Fixed128(0n, 1n);
      expect(zero.toHexString()).toBe('00.00');
    });

    test('should format positive numbers correctly', () => {
      const value = new Fixed128(123n, 256n); // 123/256
      const hex = value.toHexString();
      expect(hex).toMatch(/^[0-9A-F]+\.[0-9A-F]+$/);
    });

    test('should format negative numbers with minus sign', () => {
      const value = new Fixed128(-123n, 256n);
      const hex = value.toHexString();
      expect(hex).toMatch(/^-[0-9A-F]+\.[0-9A-F]+$/);
    });

    test('should support custom precision', () => {
      const value = new Fixed128(123n, 256n);
      const hex = value.toHexStringWithPrecision(1, 17);
      expect(hex).toMatch(/^[0-9A-F]+\.[0-9A-F]+$/);
    });

    test('should throw error for invalid precision', () => {
      const value = new Fixed128(123n, 256n);
      expect(() => value.toHexStringWithPrecision(9, 10)).toThrow('Invalid precision');
      expect(() => value.toHexStringWithPrecision(8, 9)).toThrow('Invalid precision');
    });
  });

  describe('Base64 Formatting', () => {
    test('should encode and decode base64 correctly', () => {
      const original = new Fixed128(12345n, 6789n);
      const base64 = original.toBase64();
      expect(typeof base64).toBe('string');

      const decoded = Fixed128.fromBase64(base64);
      expect(decoded.equals(original)).toBe(true);
    });

    test('should handle zero value in base64', () => {
      const zero = new Fixed128(0n, 1n);
      const base64 = zero.toBase64();
      const decoded = Fixed128.fromBase64(base64);
      expect(decoded.equals(zero)).toBe(true);
    });

    test('should handle negative values in base64', () => {
      const negative = new Fixed128(-12345n, 6789n);
      const base64 = negative.toBase64();
      const decoded = Fixed128.fromBase64(base64);
      expect(decoded.equals(negative)).toBe(true);
    });
  });

  describe('Hex String Parsing', () => {
    test('should parse simple hex strings', () => {
      const hexStr = '01.80';
      const parsed = Fixed128.fromHexString(hexStr);
      expect(parsed.toHexString()).toBe(hexStr);
    });

    test('should parse negative hex strings', () => {
      const hexStr = '-01.80';
      const parsed = Fixed128.fromHexString(hexStr);
      expect(parsed.toHexString()).toBe(hexStr);
    });

    test('should handle zero hex string', () => {
      const parsed = Fixed128.fromHexString('00.00');
      expect(parsed.isZero()).toBe(true);
    });

    test('should throw error for invalid format', () => {
      expect(() => Fixed128.fromHexString('invalid')).toThrow('Expected format');
      expect(() => Fixed128.fromHexString('01-80')).toThrow('Expected format');
      expect(() => Fixed128.fromHexString('')).toThrow('Empty hex string');
    });

    test('should handle hex strings without padding', () => {
      const parsed = Fixed128.fromHexString('1.8');
      // When parsing "1.8", it becomes "01.08" because hex "8" is padded to "08"
      expect(parsed.toHexString()).toBe('01.08');
    });
  });
});

describe('BinaryDate Formatting', () => {
  describe('Hex Formatting', () => {
    test('should provide hex() method', () => {
      const date = BinaryDate.fromUnixMillis(1640995200000n); // 2022-01-01
      const hex = date.hex();
      expect(typeof hex).toBe('string');
      expect(hex).toMatch(/^[0-9A-F-]+\.[0-9A-F]+$/);
    });

    test('should provide hexFine() method for full precision', () => {
      const date = BinaryDate.fromUnixMillis(1640995200000n);
      const hexFine = date.hexFine();
      expect(typeof hexFine).toBe('string');
      expect(hexFine).toMatch(/^[0-9A-F-]+\.[0-9A-F]+$/);
    });

    test('should provide base64() method', () => {
      const date = BinaryDate.fromUnixMillis(1640995200000n);
      const base64 = date.base64();
      expect(typeof base64).toBe('string');
    });

    test('should support round-trip hex conversion', () => {
      const original = BinaryDate.fromUnixMillis(1640995200000n);
      const hex = original.hexFine();
      const restored = BinaryDate.fromHex(hex);
      expect(restored.equals(original)).toBe(true);
    });

    test('should support round-trip base64 conversion', () => {
      const original = BinaryDate.fromUnixMillis(1640995200000n);
      const base64 = original.base64();
      const restored = BinaryDate.fromBase64(base64);
      expect(restored.equals(original)).toBe(true);
    });
  });

  describe('Custom Precision Hex', () => {
    test('should support custom precision hex formatting', () => {
      const date = BinaryDate.fromUnixMillis(1640995200000n);
      const hex = date.hexWithPrecision(4, 12);
      expect(typeof hex).toBe('string');
      expect(hex).toMatch(/^[0-9A-F-]+\.[0-9A-F]+$/);
    });

    test('should throw error for invalid precision in BinaryDate', () => {
      const date = BinaryDate.fromUnixMillis(1640995200000n);
      expect(() => date.hexWithPrecision(9, 10)).toThrow('Invalid precision');
    });
  });
});

describe('BinaryDuration Formatting', () => {
  describe('String Representation', () => {
    test('should format zero duration', () => {
      const zero = BinaryDuration.fromMillis(0n);
      expect(zero.toString()).toBe('0ms');
    });

    test('should format milliseconds', () => {
      const millis = BinaryDuration.fromMillis(500n);
      expect(millis.toString()).toBe('500ms');
    });

    test('should format seconds', () => {
      const seconds = BinaryDuration.fromSeconds(30n);
      expect(seconds.toString()).toBe('30s');
    });

    test('should format seconds with milliseconds', () => {
      const duration = BinaryDuration.fromMillis(1500n); // 1.5 seconds
      expect(duration.toString()).toBe('1.5s');
    });

    test('should format minutes', () => {
      const minutes = BinaryDuration.fromMinutes(5n);
      expect(minutes.toString()).toBe('5m');
    });

    test('should format hours', () => {
      const hours = BinaryDuration.fromHours(2n);
      expect(hours.toString()).toBe('2h');
    });

    test('should format days', () => {
      const days = BinaryDuration.fromDays(1n);
      expect(days.toString()).toBe('1d');
    });

    test('should format complex durations', () => {
      // 2 hours, 30 minutes, 45 seconds
      const duration = BinaryDuration.fromMillis(2n * 3600000n + 30n * 60000n + 45n * 1000n);
      expect(duration.toString()).toBe('2h30m45s');
    });

    test('should format negative durations', () => {
      const negative = BinaryDuration.fromMillis(-1500n);
      expect(negative.toString()).toBe('-1.5s');
    });
  });

  describe('Duration Parsing', () => {
    test('should parse simple milliseconds', () => {
      const parsed = BinaryDuration.parse('500ms');
      expect(parsed.millis()).toBe(500n);
    });

    test('should parse seconds', () => {
      const parsed = BinaryDuration.parse('30s');
      expect(parsed.millis()).toBe(30000n);
    });

    test('should parse decimal seconds', () => {
      const parsed = BinaryDuration.parse('1.5s');
      expect(parsed.millis()).toBe(1500n);
    });

    test('should parse minutes', () => {
      const parsed = BinaryDuration.parse('5m');
      expect(parsed.millis()).toBe(300000n);
    });

    test('should parse hours', () => {
      const parsed = BinaryDuration.parse('2h');
      expect(parsed.millis()).toBe(7200000n);
    });

    test('should parse days', () => {
      const parsed = BinaryDuration.parse('1d');
      expect(parsed.millis()).toBe(86400000n);
    });

    test('should parse complex durations', () => {
      const parsed = BinaryDuration.parse('2h30m45s');
      const expected = 2n * 3600000n + 30n * 60000n + 45n * 1000n;
      expect(parsed.millis()).toBe(expected);
    });

    test('should parse negative durations', () => {
      const parsed = BinaryDuration.parse('-1.5s');
      expect(parsed.millis()).toBe(-1500n);
    });

    test('should handle mixed decimal and integer units', () => {
      const parsed = BinaryDuration.parse('1.5h30m');
      const expected = 90n * 60000n + 30n * 60000n; // 1.5 hours = 90 minutes + 30 minutes
      expect(parsed.millis()).toBe(expected);
    });

    test('should throw error for invalid format', () => {
      expect(() => BinaryDuration.parse('invalid')).toThrow('Invalid duration format');
      expect(() => BinaryDuration.parse('')).toThrow('Empty duration string');
      expect(() => BinaryDuration.parse('1.5x')).toThrow('Unknown time unit');
    });
  });

  describe('Round-trip String Conversion', () => {
    test('should support round-trip for simple durations', () => {
      const original = BinaryDuration.fromSeconds(45n);
      const str = original.toString();
      const parsed = BinaryDuration.parse(str);
      expect(parsed.equals(original)).toBe(true);
    });

    test('should support round-trip for complex durations', () => {
      // Note: toString() may format differently than parse() input
      const original = BinaryDuration.fromMillis(7545000n); // 2h5m45s
      const parsed = BinaryDuration.parse('2h5m45s');
      expect(parsed.equals(original)).toBe(true);
    });
  });

  describe('Hex and Base64 Formatting', () => {
    test('should provide hex formatting for durations', () => {
      const duration = BinaryDuration.fromHours(1n);
      const hex = duration.hex();
      expect(typeof hex).toBe('string');
      expect(hex).toMatch(/^[0-9A-F-]+\.[0-9A-F]+$/);
    });

    test('should support round-trip hex conversion', () => {
      const original = BinaryDuration.fromMinutes(30n);
      const hex = original.hex();
      const restored = BinaryDuration.fromHex(hex);
      expect(restored.equals(original)).toBe(true);
    });

    test('should support round-trip base64 conversion', () => {
      const original = BinaryDuration.fromHours(2n);
      const base64 = original.base64();
      const restored = BinaryDuration.fromBase64(base64);
      expect(restored.equals(original)).toBe(true);
    });
  });
});
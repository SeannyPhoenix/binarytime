import { Fixed128 } from '../lib/Fixed128';

describe('Fixed128', () => {
  describe('constructor', () => {
    it('should create Fixed128 from positive fraction', () => {
      const f = new Fixed128(1n, 2n);
      expect(f).toBeInstanceOf(Fixed128);
      expect(f.isZero()).toBe(false);
      expect(f.isNeg()).toBe(false);
    });

    it('should create Fixed128 from negative fraction', () => {
      const f = new Fixed128(-1n, 2n);
      expect(f.isNeg()).toBe(true);
    });

    it('should throw on division by zero', () => {
      expect(() => new Fixed128(1n, 0n)).toThrow('Division by zero');
    });
  });

  describe('static constants', () => {
    it('should have ZERO constant', () => {
      expect(Fixed128.ZERO.isZero()).toBe(true);
    });

    it('should have ONE constant', () => {
      expect(Fixed128.ONE.isZero()).toBe(false);
      expect(Fixed128.ONE.sign()).toBe(1);
    });
  });

  describe('arithmetic operations', () => {
    it('should add two Fixed128 numbers', () => {
      const a = new Fixed128(1n, 4n); // 0.25
      const b = new Fixed128(1n, 4n); // 0.25
      const result = a.add(b); // 0.5

      expect(result).toBeInstanceOf(Fixed128);
      expect(result.equals(new Fixed128(1n, 2n))).toBe(true);
    });

    it('should subtract two Fixed128 numbers', () => {
      const a = new Fixed128(3n, 4n); // 0.75
      const b = new Fixed128(1n, 4n); // 0.25
      const result = a.sub(b); // 0.5

      expect(result.equals(new Fixed128(1n, 2n))).toBe(true);
    });

    it('should multiply two Fixed128 numbers', () => {
      const a = new Fixed128(1n, 2n); // 0.5
      const b = new Fixed128(1n, 2n); // 0.5
      const result = a.mul(b); // 0.25

      expect(result).toBeInstanceOf(Fixed128);
    });

    it('should divide two Fixed128 numbers', () => {
      const a = new Fixed128(1n, 2n); // 0.5
      const b = new Fixed128(1n, 4n); // 0.25
      const result = a.quo(b); // 2.0

      expect(result).toBeInstanceOf(Fixed128);
    });

    it('should throw on division by zero in quo', () => {
      const a = new Fixed128(1n, 2n);
      expect(() => a.quo(Fixed128.ZERO)).toThrow('Division by zero');
    });
  });

  describe('comparison operations', () => {
    it('should compare Fixed128 numbers correctly', () => {
      const a = new Fixed128(1n, 4n); // 0.25
      const b = new Fixed128(1n, 2n); // 0.5
      const c = new Fixed128(1n, 4n); // 0.25

      expect(a.cmp(b)).toBe(-1); // a < b
      expect(b.cmp(a)).toBe(1);  // b > a
      expect(a.cmp(c)).toBe(0);  // a == c
    });

    it('should check equality correctly', () => {
      const a = new Fixed128(1n, 4n);
      const b = new Fixed128(1n, 4n);
      const c = new Fixed128(1n, 2n);

      expect(a.equals(b)).toBe(true);
      expect(a.equals(c)).toBe(false);
    });
  });

  describe('sign operations', () => {
    it('should return correct sign', () => {
      const positive = new Fixed128(1n, 2n);
      const negative = new Fixed128(-1n, 2n);
      const zero = Fixed128.ZERO;

      expect(positive.sign()).toBe(1);
      expect(negative.sign()).toBe(-1);
      expect(zero.sign()).toBe(0);
    });

    it('should detect negative numbers', () => {
      const positive = new Fixed128(1n, 2n);
      const negative = new Fixed128(-1n, 2n);

      expect(positive.isNeg()).toBe(false);
      expect(negative.isNeg()).toBe(true);
    });
  });

  describe('immutability', () => {
    it('should not modify original instances during operations', () => {
      const original = new Fixed128(1n, 2n);
      const originalValue = original.getValue();

      original.add(new Fixed128(1n, 4n));

      expect(original.getValue()).toBe(originalValue);
    });

    it('should return same instance from copy', () => {
      const original = new Fixed128(1n, 2n);
      const copy = original.copy();

      expect(copy).toBe(original); // Same reference since immutable
    });
  });

  describe('conversion methods', () => {
    it('should convert to bytes', () => {
      const f = new Fixed128(1n, 2n);
      const bytes = f.bytes();

      expect(bytes).toBeInstanceOf(Uint8Array);
      expect(bytes.length).toBe(16);
    });

    it('should get hi/lo components', () => {
      const f = new Fixed128(1n, 2n);
      const [hi, lo] = f.hiLo();

      expect(typeof hi).toBe('bigint');
      expect(typeof lo).toBe('bigint');
    });

    it('should have string representation', () => {
      const f = new Fixed128(1n, 2n);
      const str = f.toString();

      expect(typeof str).toBe('string');
      expect(str).toContain('Fixed128');
    });
  });
});
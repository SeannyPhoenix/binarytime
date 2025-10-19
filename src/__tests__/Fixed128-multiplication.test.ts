import { Fixed128 } from '../lib/Fixed128';

describe('Fixed128 multiplication optimizations', () => {
  it('should return zero when multiplying by zero', () => {
    const f = new Fixed128(5n, 3n); // Some non-zero value
    const result = f.mulBigInt(0n);

    expect(result).toBe(0n);
  });

  it('should handle negative numbers multiplied by zero', () => {
    const f = new Fixed128(-7n, 2n); // Some negative value
    const result = f.mulBigInt(0n);

    expect(result).toBe(0n);
  });

  it('should handle zero Fixed128 multiplied by any number', () => {
    const f = Fixed128.ZERO;
    const result = f.mulBigInt(1000n);

    expect(result).toBe(0n);
  });
});
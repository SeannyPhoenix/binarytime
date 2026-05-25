import * as z from "zod";

const DAY_MS = 86_400_000n;

// Added to a Binary Time calculated from the epoch
// to ensure no negatives and to allow for dates
// far in the past and future without overflow issues.
// Approximately 24.08 billion years in BinaryTime
const BINARY_TIME_OFFSET = (1n << 43n) << 64n;

// Exactly 16 hex digits, a period, and 16 hex digits
// Represents 128 bits, or 16 bytes
// This representation is 33 bytes long
const BINARY_TIME_REGEX = /^[0-9a-f]{16}\.[0-9a-f]{16}$/;

export type BinaryTime = bigint;

export function binaryTimeFromString(s: string): BinaryTime {
  if (!BINARY_TIME_REGEX.test(s)) {
    throw new Error("Invalid binary time format");
  }

  const [hi, lo] = s.split(".").map((part) => BigInt(`0x${part}`));
  return (hi << 64n) | lo;
}

export function binaryTimeToString(b: BinaryTime): string {
  const whole = b >> 64n;
  const frac = b & 0xffffffffffffffffn;

  const hi = whole.toString(16).padStart(16, "0");
  const lo = frac.toString(16).padStart(16, "0");
  return `${hi}.${lo}`;
}

const zBinaryTimeString = z
  .string()
  .regex(BINARY_TIME_REGEX, { error: "Invalid binary time format" });

export const zBinaryTime = z.codec(zBinaryTimeString, z.bigint(), {
  encode: binaryTimeToString,
  decode: binaryTimeFromString,
});

export function binaryTimeNow(): BinaryTime {
  return binaryTimeFromMs(BigInt(Date.now()));
}

export function binaryTimeFromDate(date: Date): BinaryTime {
  return binaryTimeFromMs(BigInt(date.getTime()));
}

export function binaryTimeFromMs(ms: bigint): BinaryTime {
  const whole = ms / DAY_MS;
  const remainder = ms % DAY_MS;
  const frac = (remainder << 64n) / DAY_MS;

  const deltaFromEpoch = (whole << 64n) | frac;
  return deltaFromEpoch + BINARY_TIME_OFFSET;
}

export function binaryTimeToMs(binarytime: BinaryTime): number {
  const deltaFromEpoch = binarytime - BINARY_TIME_OFFSET;
  const whole = deltaFromEpoch >> 64n;
  const frac = deltaFromEpoch & 0xffffffffffffffffn;

  return Number(whole * DAY_MS + ((frac * DAY_MS + (1n << 63n)) >> 64n));
}

export function binaryTimeToDate(binarytime: BinaryTime): Date {
  return new Date(binaryTimeToMs(binarytime));
}

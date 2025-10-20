// Export core classes
export { Fixed128 } from './lib/Fixed128';
export { BinaryDate } from './lib/BinaryDate';
export { BinaryDuration } from './lib/BinaryDuration';

// Export constants
export {
  DAY_MILLISECONDS,
  DAY_SECONDS,
  DAY_MINUTES,
  DAY_HOURS,
} from './lib/constants';

// Convenience functions for BinaryDate and BinaryDuration arithmetic
import { BinaryDate } from './lib/BinaryDate';
import { BinaryDuration } from './lib/BinaryDuration';

/**
 * Add a BinaryDuration to a BinaryDate
 * @param date - the date to add to
 * @param duration - the duration to add
 * @returns new BinaryDate with the duration added
 */
export function addDuration(date: BinaryDate, duration: BinaryDuration): BinaryDate {
  return date.addFixed128Duration(duration.fixed128());
}

/**
 * Subtract a BinaryDuration from a BinaryDate
 * @param date - the date to subtract from
 * @param duration - the duration to subtract
 * @returns new BinaryDate with the duration subtracted
 */
export function subDuration(date: BinaryDate, duration: BinaryDuration): BinaryDate {
  return date.subFixed128Duration(duration.fixed128());
}

/**
 * Calculate the duration between two BinaryDates
 * @param from - the starting date
 * @param to - the ending date
 * @returns BinaryDuration representing the time span (positive if to > from)
 */
export function durationBetween(from: BinaryDate, to: BinaryDate): BinaryDuration {
  const diff = from.durationUntilFixed128(to);
  return new BinaryDuration(diff);
}

/**
 * Calculate the duration since a BinaryDate
 * @param from - the starting date
 * @param to - the ending date (usually later)
 * @returns BinaryDuration representing the time span since from
 */
export function durationSince(from: BinaryDate, to: BinaryDate): BinaryDuration {
  const diff = to.durationSinceFixed128(from);
  return new BinaryDuration(diff);
}

// Version information
export const VERSION = '0.1.0';
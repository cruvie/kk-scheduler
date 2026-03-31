/**
 * Time utility for handling default time values
 * Matches Go's DefaultTime: time.Date(1999, 1, 12, 0, 0, 0, 0, time.UTC)
 */

// Default time: 1999-01-12T00:00:00.000Z
export const DEFAULT_TIME = new Date(Date.UTC(1999, 0, 12, 0, 0, 0, 0));

// Timestamp type from protobuf
type Timestamp = { seconds: bigint | number } | null | undefined;

/**
 * Check if a timestamp is the default time
 * @param ts - Timestamp with seconds field
 * @returns true if the timestamp matches the default time
 */
export function isDefaultTime(ts: Timestamp): boolean {
  if (!ts?.seconds) return true;
  const date = new Date(Number(ts.seconds) * 1000);
  if (isNaN(date.getTime())) return true;
  return date.getTime() === DEFAULT_TIME.getTime();
}

/**
 * Format a timestamp for display, returns "-" if it's the default time
 * @param ts - Timestamp with seconds field
 * @returns Formatted date string or "-" for default time
 */
export function displayTime(ts: Timestamp): string {
  if (isDefaultTime(ts)) {
    return '-';
  }
  const date = new Date(Number(ts!.seconds) * 1000);
  // Format: YYYY-MM-DD HH:mm:ss
  return date.toLocaleString('sv-SE', {
    timeZone: 'UTC',
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  }).replace('T', ' ');
}
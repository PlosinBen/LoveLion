// Dates are stored as naive timestamps (no timezone). The Go backend
// always serialises time.Time with a Z suffix, so we strip it before
// parsing to prevent JavaScript from applying a UTC→local shift.
export function parseNaiveDate(dateStr: string): Date {
  return new Date(dateStr.replace(/Z$/, ''))
}

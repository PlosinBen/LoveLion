// Dates are stored as naive timestamps (no timezone). The Go backend
// always serialises time.Time with a Z suffix, so we strip it before
// parsing to prevent JavaScript from applying a UTC→local shift.
export function parseNaiveDate(dateStr: string): Date {
  return new Date(dateStr.replace(/Z$/, ''))
}

// Mirror of parseNaiveDate for the write path: serialise a Date using its
// local components with a fake Z suffix so the backend stores exactly what
// the user saw on screen (e.g. 09:00 Taipei stays 09:00, never 01:00 UTC).
// Do NOT use Date#toISOString — that converts to real UTC and shifts the time.
export function toLocalISOString(d: Date): string {
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}Z`
}

// Today as YYYY-MM-DD in the browser's local timezone. Used for date-only
// form defaults. Do NOT use `new Date().toISOString().slice(0, 10)` — that
// uses UTC and shifts to the previous day during Taipei midnight–08:00.
export function todayLocalDate(): string {
  const d = new Date()
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
}

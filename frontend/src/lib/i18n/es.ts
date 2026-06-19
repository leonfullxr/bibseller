// Spanish catalogue. Filled in M8.2 (#46); for the M8.1 foundation it is
// intentionally empty, so `/es` renders the English fallback while the routing,
// switcher, formatting and hreflang machinery is exercised end to end. Any key
// missing here falls back to en.ts via t().
import type { MessageKey } from './en';

export const es: Partial<Record<MessageKey, string>> = {};

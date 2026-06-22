// Fails CI when a .svelte file carries a hardcoded UI string instead of a t()
// key (M8.1 / D17). Deliberately a heuristic, not a parser (D14: no new deps):
// it blanks <script>/<style>/comments and {...} expressions - preserving line
// numbers - then flags leftover literal text between tags and literal values of
// the user-facing attributes below. Dynamic data ({race.sport}, {form.error})
// and t(...) output are expressions, so they are never flagged. It also flags an
// href/action that calls resolve() without a link() wrapper, which would drop
// the /es prefix on the Spanish site.
//
// Scope is .svelte UI only. Server-side / API error strings are keyed too (#49)
// but enforced by types + review, not this heuristic; this guard does not scan them.
//
// Run: node scripts/check-i18n.mjs  (also wired into `npm run lint`).
import { readdirSync, readFileSync, statSync } from 'node:fs';
import { join, relative } from 'node:path';
import { fileURLToPath } from 'node:url';

// fileURLToPath (not .pathname) so the path is correct on Windows too.
const ROOT = fileURLToPath(new URL('../src', import.meta.url));
const ATTRS = ['placeholder', 'aria-label', 'title', 'alt'];

// Literal text that is intentionally not translated: the brand wordmark and the
// `make dev` command shown verbatim. Matched against the trimmed text segment.
const ALLOW = new Set(['bib', 'seller', 'make dev']);

// Replace every non-newline char with a space, so blanked regions keep their
// line/column footprint and reported line numbers stay accurate.
const blank = (m) => m.replace(/[^\n]/g, ' ');
const lineOf = (src, index) => src.slice(0, index).split('\n').length;

/** Recursively collect .svelte files under dir. */
function svelteFiles(dir) {
	const out = [];
	for (const entry of readdirSync(dir)) {
		const full = join(dir, entry);
		if (statSync(full).isDirectory()) out.push(...svelteFiles(full));
		else if (entry.endsWith('.svelte')) out.push(full);
	}
	return out;
}

function violations(src) {
	const found = [];
	// 1. Drop script/style/comments (keep line numbers).
	const masked = src
		.replace(/<script[\s\S]*?<\/script>/g, blank)
		.replace(/<style[\s\S]*?<\/style>/g, blank)
		.replace(/<!--[\s\S]*?-->/g, blank);

	// 2. Literal user-facing attribute values (t() output is {..}, never "..").
	const attrRe = new RegExp(`\\b(${ATTRS.join('|')})\\s*=\\s*"([^"]*)"`, 'g');
	for (let m; (m = attrRe.exec(masked)); ) {
		const value = m[2].trim();
		if (/[A-Za-z]{2}/.test(value) && !ALLOW.has(value)) {
			found.push({ line: lineOf(masked, m.index), text: `${m[1]}="${value}"` });
		}
	}

	// 3. Internal nav must keep the locale prefix: an href/action that calls
	// resolve() without a link() wrapper drops /es on the Spanish site. Matches
	// both unquoted `={...}` and quoted `="{...}?{qs}"` forms, reading up to the
	// first } or " - enough to spot a bare resolve.
	const navRe = /\b(href|action)=(?:\{|")([^"}]*)/g;
	for (let m; (m = navRe.exec(masked)); ) {
		if (/\bresolve\(/.test(m[2]) && !/\blink\(/.test(m[2])) {
			found.push({
				line: lineOf(masked, m.index),
				text: `${m[1]}={resolve(...)} not wrapped in link()`
			});
		}
	}

	// 4. Literal text between tags: blank {..} expressions (innermost-out, so
	// nested {a, {b}} fully clears) and tags, then scan what is left.
	let text = masked;
	let prev;
	do {
		prev = text;
		text = text.replace(/\{[^{}]*\}/g, blank);
	} while (text !== prev);
	text = text.replace(/<[^>]*>/g, blank);

	const wordRe = /[^\n<>]*[A-Za-z]{2}[^\n<>]*/g;
	for (let m; (m = wordRe.exec(text)); ) {
		// A run of 2+ spaces is a blanked tag/expression, i.e. a boundary between
		// separate text nodes (e.g. brand `bib<span>seller</span>`) - split on it
		// so adjacent allowlisted nodes are judged individually, not merged.
		for (const part of m[0].split(/\s{2,}/)) {
			const seg = part.trim();
			if (seg && /[A-Za-z]{2}/.test(seg) && !ALLOW.has(seg)) {
				found.push({ line: lineOf(text, m.index), text: seg });
			}
		}
	}
	return found;
}

let total = 0;
for (const file of svelteFiles(ROOT)) {
	const hits = violations(readFileSync(file, 'utf8'));
	for (const h of hits) {
		console.error(
			`${relative(ROOT, file)}:${h.line}  hardcoded UI string: ${JSON.stringify(h.text)}`
		);
		total++;
	}
}

if (total > 0) {
	console.error(
		`\n${total} hardcoded UI string(s) found. Move them into $lib/i18n/en.ts and use t().`
	);
	process.exit(1);
}
console.log('i18n: no hardcoded UI strings in .svelte files.');

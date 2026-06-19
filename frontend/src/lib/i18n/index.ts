// Component-facing i18n: the bound { locale, t, link } the active locale
// produces, shared down the tree via Svelte context (per-render-tree, so
// SSR-safe - no module-level locale that could bleed across concurrent
// requests). The +layout builds the instance with runes so locale stays
// reactive; this module only owns the type + the context key.
//
// `link` prefixes an already-resolved path with `/es` when Spanish is active.
// Call it as `link(resolve('/races'))`: resolve keeps full route/param checking,
// and because link returns `ResolvedPathname`, `<a href={link(resolve(...))}>`
// is exactly what svelte/no-navigation-without-resolve accepts (it checks the
// value's type, not a literal resolve() call).
import { getContext, setContext } from 'svelte';
import type { ResolvedPathname } from '$app/types';
import type { Locale } from './locale';
import type { Pluralizer, Translator } from './messages';

export interface I18n {
	locale: Locale;
	t: Translator;
	plural: Pluralizer;
	link: (path: ResolvedPathname) => ResolvedPathname;
}

const I18N_KEY = Symbol('i18n');

export function setI18n(i18n: I18n): void {
	setContext(I18N_KEY, i18n);
}

export function getI18n(): I18n {
	return getContext(I18N_KEY);
}

export * from './locale';
export { createPlural, createTranslator, type Pluralizer, type Translator } from './messages';
export type { MessageKey } from './en';

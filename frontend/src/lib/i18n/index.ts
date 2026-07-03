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

// Components rendered under the root layout always have the context. getI18n
// throws if it is missing (a "forgot setI18n" bug) rather than handing back
// undefined that crashes obscurely on the first t(). The one place that can
// render without it - +error.svelte, above the layout - uses tryGetI18n.
export function getI18n(): I18n {
	const i18n = getContext<I18n | undefined>(I18N_KEY);
	if (!i18n) throw new Error('getI18n() called outside an i18n context (root layout missing?)');
	return i18n;
}

export function tryGetI18n(): I18n | undefined {
	return getContext<I18n | undefined>(I18N_KEY);
}

export * from './locale';
export {
	createPlural,
	createTranslator,
	listingStatusLabel,
	sportLabel,
	type Pluralizer,
	type Translator
} from './messages';
export type { MessageKey } from './en';

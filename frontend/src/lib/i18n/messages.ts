// The translator: framework-free so unit tests import it without a Kit runtime.
// createTranslator binds a catalogue to the active locale; t(key, params) looks
// the key up, falls back to English when a locale has not translated it, and
// fills {placeholder} params. Missing-everywhere keys are impossible by typing
// (MessageKey = keyof en), so there is no runtime "key not found" branch.
import type { Locale } from './locale';
import { en, type MessageKey } from './en';
import { es } from './es';

const messages: Record<Locale, Partial<Record<MessageKey, string>>> = { en, es };

type Params = Record<string, string | number>;
export type Translator = (key: MessageKey, params?: Params) => string;
export type Pluralizer = (base: string, n: number, params?: Params) => string;

function interpolate(template: string, params?: Params): string {
	if (!params) return template;
	return template.replace(/\{(\w+)\}/g, (_, name) =>
		name in params ? String(params[name]) : `{${name}}`
	);
}

export function createTranslator(locale: Locale): Translator {
	const dict = messages[locale];
	return (key, params) => interpolate(dict[key] ?? en[key], params);
}

// createPlural resolves `${base}.${category}`, where category is the CLDR plural
// form for n in this locale via Intl.PluralRules (stdlib, no dependency). en/es
// only need one|other; a language with more forms (Polish, Arabic) just needs
// the extra `${base}.few` / `.many` keys added - no code change. Falls back to
// `${base}.other`, then English, so a missing form never renders undefined.
export function createPlural(locale: Locale): Pluralizer {
	const dict = messages[locale];
	const rules = new Intl.PluralRules(locale);
	return (base, n, params) => {
		const key = `${base}.${rules.select(n)}` as MessageKey;
		const other = `${base}.other` as MessageKey;
		return interpolate(dict[key] ?? en[key] ?? dict[other] ?? en[other] ?? '', { n, ...params });
	};
}

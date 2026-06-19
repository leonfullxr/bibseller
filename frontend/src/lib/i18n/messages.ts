// The translator: framework-free so unit tests import it without a Kit runtime.
// createTranslator binds a catalogue to the active locale; t(key, params) looks
// the key up, falls back to English when a locale has not translated it, and
// fills {placeholder} params. Missing-everywhere keys are impossible by typing
// (MessageKey = keyof en), so there is no runtime "key not found" branch.
import type { Locale } from './locale';
import { en, type MessageKey } from './en';
import { es } from './es';

const messages: Record<Locale, Partial<Record<MessageKey, string>>> = { en, es };

export type Translator = (key: MessageKey, params?: Record<string, string | number>) => string;

export function createTranslator(locale: Locale): Translator {
	const dict = messages[locale];
	return (key, params) => {
		const template = dict[key] ?? en[key];
		if (!params) return template;
		return template.replace(/\{(\w+)\}/g, (_, name) =>
			name in params ? String(params[name]) : `{${name}}`
		);
	};
}

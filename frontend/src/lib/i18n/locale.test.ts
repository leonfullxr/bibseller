import { describe, expect, it } from 'vitest';
import {
	detectFromAcceptLanguage,
	isBot,
	localeFromPath,
	pathForLocale,
	stripLocale,
	suggestsSpanish
} from './locale';
import { createPlural, createTranslator } from './messages';
import { en } from './en';
import { transferPolicies } from '$lib/policy';

describe('localeFromPath / stripLocale / pathForLocale', () => {
	it('reads the locale from the path prefix', () => {
		expect(localeFromPath('/')).toBe('en');
		expect(localeFromPath('/races')).toBe('en');
		expect(localeFromPath('/es')).toBe('es');
		expect(localeFromPath('/es/races')).toBe('es');
		expect(localeFromPath('/espana')).toBe('en'); // not a locale prefix
	});

	it('strips and re-applies the prefix as inverses', () => {
		for (const path of ['/', '/races', '/listings/abc']) {
			expect(stripLocale(pathForLocale('es', path))).toBe(path);
			expect(stripLocale(pathForLocale('en', path))).toBe(path);
		}
		expect(pathForLocale('es', '/')).toBe('/es');
		expect(pathForLocale('es', '/races')).toBe('/es/races');
		expect(pathForLocale('en', '/races')).toBe('/races');
		expect(stripLocale('/es')).toBe('/');
	});
});

describe('detectFromAcceptLanguage', () => {
	it('falls back to en when missing or unsupported', () => {
		expect(detectFromAcceptLanguage(null)).toBe('en');
		expect(detectFromAcceptLanguage('fr-FR,fr;q=0.9')).toBe('en');
	});

	it('honours quality weights over header order', () => {
		expect(detectFromAcceptLanguage('es-ES,es;q=0.9,en;q=0.8')).toBe('es');
		expect(detectFromAcceptLanguage('en-US,en;q=0.9,es;q=0.8')).toBe('en');
		expect(detectFromAcceptLanguage('en;q=0.7,es;q=0.95')).toBe('es');
	});

	it('coerces a malformed q to 0 (deterministic, never NaN)', () => {
		// es has a broken q -> 0, so the well-formed en wins deterministically.
		expect(detectFromAcceptLanguage('es;q=abc,en;q=0.5')).toBe('en');
	});
});

describe('suggestsSpanish', () => {
	it('uses geo country when present (case-insensitive), ignoring the browser', () => {
		expect(suggestsSpanish('ES', null)).toBe(true);
		expect(suggestsSpanish('es', null)).toBe(true);
		expect(suggestsSpanish('FR', 'es-ES,es;q=0.9')).toBe(false); // geo wins
	});

	it('falls back to Accept-Language when geo is unknown', () => {
		expect(suggestsSpanish(null, 'es-ES,es;q=0.9')).toBe(true);
		expect(suggestsSpanish('', 'en-US,en')).toBe(false);
		expect(suggestsSpanish(null, null)).toBe(false);
	});
});

describe('isBot', () => {
	it('matches common crawlers, not real browsers', () => {
		expect(isBot('Googlebot/2.1')).toBe(true);
		expect(isBot('facebookexternalhit/1.1')).toBe(true);
		expect(isBot(null)).toBe(false);
		expect(isBot('Mozilla/5.0 (Macintosh) Safari/605')).toBe(false);
	});
});

describe('createTranslator', () => {
	it('returns the English string and falls back when es is unfilled', () => {
		const tEn = createTranslator('en');
		const tEs = createTranslator('es');
		expect(tEn('nav.races')).toBe('Races');
		// es.ts is empty in M8.1, so every key falls back to English.
		expect(tEs('nav.races')).toBe('Races');
	});

	it('interpolates {placeholder} params', () => {
		const t = createTranslator('en');
		expect(t('listingCard.listedBy', { name: 'Ana' })).toBe('Listed by Ana');
		expect(t('raceCard.bibs.other', { n: 3 })).toBe('3 bibs listed');
	});
});

describe('createPlural', () => {
	it('selects the CLDR plural form for n and fills {n}', () => {
		const p = createPlural('en');
		expect(p('raceCard.bibs', 1)).toBe('1 bib listed');
		expect(p('raceCard.bibs', 0)).toBe('0 bibs listed');
		expect(p('raceCard.bibs', 5)).toBe('5 bibs listed');
	});
});

describe('policy dictionary coverage', () => {
	it('has a label and a disclaimer title/body for every transfer policy', () => {
		// The words policy.ts used to own now live here; guard against a missing mode.
		for (const p of transferPolicies) {
			expect(en[`policy.label.${p}`]).toBeTruthy();
			expect(en[`policy.disclaimer.${p}.title`]).toBeTruthy();
			expect(en[`policy.disclaimer.${p}.body`]).toBeTruthy();
		}
	});
});

import { describe, expect, it } from 'vitest';
import { apiUrl } from './url';

describe('apiUrl', () => {
	it('joins base and path with exactly one slash', () => {
		expect(apiUrl('http://localhost:8080', '/api/healthz')).toBe(
			'http://localhost:8080/api/healthz'
		);
		expect(apiUrl('http://localhost:8080/', 'api/healthz')).toBe(
			'http://localhost:8080/api/healthz'
		);
		expect(apiUrl('http://localhost:8080/', '/api/healthz')).toBe(
			'http://localhost:8080/api/healthz'
		);
	});

	it('keeps query strings intact', () => {
		expect(apiUrl('https://api.example.com', '/api/v1/races?country=ES')).toBe(
			'https://api.example.com/api/v1/races?country=ES'
		);
	});
});

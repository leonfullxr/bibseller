import { describe, expect, it } from 'vitest';
import { clientIPHeader } from './clientip';

describe('clientIPHeader', () => {
	it('forwards the inbound CF-Connecting-IP verbatim', () => {
		const request = new Request('http://localhost/login', {
			headers: { 'CF-Connecting-IP': '203.0.113.7' }
		});
		expect(clientIPHeader(request)).toEqual({ 'CF-Connecting-IP': '203.0.113.7' });
	});

	it('sends nothing when the header is absent (dev)', () => {
		const request = new Request('http://localhost/login');
		expect(clientIPHeader(request)).toEqual({});
	});
});

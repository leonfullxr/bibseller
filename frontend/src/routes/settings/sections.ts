// Section routing for /settings: the active pane comes from ?section=, so
// deep links, SSR, and back/forward work with zero client-side nav state.
export const sections = ['profile', 'security', 'account'] as const;
export type Section = (typeof sections)[number];

export function activeSection(param: string | null): Section {
	return (sections as readonly string[]).includes(param ?? '') ? (param as Section) : 'profile';
}

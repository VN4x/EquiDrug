import { writable } from 'svelte/store';
import type { Locale } from '$lib/i18n';

const stored =
	typeof localStorage !== 'undefined'
		? (localStorage.getItem('equidrug-locale') as Locale | null)
		: null;

export const locale = writable<Locale>(stored ?? 'en');

locale.subscribe((value) => {
	if (typeof localStorage !== 'undefined') {
		localStorage.setItem('equidrug-locale', value);
	}
});

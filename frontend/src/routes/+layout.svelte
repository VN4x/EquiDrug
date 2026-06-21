<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { page } from '$app/stores';
	import { locale } from '$lib/stores/locale';
	import { locales, t, type Locale } from '$lib/i18n';

	let { children } = $props();

	const navItems = [
		{ href: '/', icon: '🔍', key: 'nav.lookup' },
		{ href: '/locker', icon: '💊', key: 'nav.locker' },
		{ href: '/planner', icon: '✈️', key: 'nav.planner' },
		{ href: '/wiki', icon: '📖', key: 'nav.wiki' },
		{ href: '/avoid', icon: '⚠️', key: 'nav.avoid' }
	];

	function isActive(href: string, pathname: string) {
		if (href === '/') return pathname === '/';
		return pathname.startsWith(href);
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover" />
	<meta name="theme-color" content="#0d6e6e" />
	<title>EquiDrug</title>
</svelte:head>

<div class="app-shell">
	<header class="app-header">
		<div class="brand">
			<h1>{t($locale, 'app.name')}</h1>
			<p>{t($locale, 'app.tagline')}</p>
		</div>
		<select
			class="locale-select"
			value={$locale}
			onchange={(e) => locale.set((e.currentTarget as HTMLSelectElement).value as Locale)}
			aria-label="Language"
		>
			{#each locales as loc}
				<option value={loc.code}>{loc.label}</option>
			{/each}
		</select>
	</header>

	<main class="app-main">
		{@render children()}
	</main>

	<nav class="bottom-nav" aria-label="Main">
		{#each navItems as item}
			<a href={item.href} class:active={isActive(item.href, $page.url.pathname)}>
				<span class="icon">{item.icon}</span>
				<span>{t($locale, item.key)}</span>
			</a>
		{/each}
	</nav>
</div>

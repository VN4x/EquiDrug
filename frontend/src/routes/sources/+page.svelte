<script lang="ts">
	import { onMount } from 'svelte';
	import DestinationSources from '$lib/components/DestinationSources.svelte';
	import {
		countries,
		getDestinationSources,
		listSourceCountries,
		type DestinationSourcesResponse
	} from '$lib/api/client';

	const specialRegions = [
		{ code: 'CN-XZ', name: 'Tibet (TAR)', parent: 'CN' },
		{ code: 'EU', name: 'European Union (fallback)', parent: '' },
		{ code: 'DEFAULT', name: 'Global fallback', parent: '' }
	];

	let countryCode = $state('VN');
	let regionCode = $state('');
	let category = $state('drug');
	let data = $state<DestinationSourcesResponse | null>(null);
	let indexed = $state<{ code: string; name: string; google_tier: string }[]>([]);
	let loading = $state(false);

	onMount(async () => {
		try {
			const res = await listSourceCountries();
			indexed = res.countries;
		} catch {
			indexed = countries.map((c) => ({ code: c.code, name: c.name, google_tier: 'secondary' }));
		}
		load();
	});

	async function load() {
		loading = true;
		try {
			data = await getDestinationSources(
				regionCode ? '' : countryCode,
				category,
				regionCode || undefined
			);
		} catch {
			data = null;
		} finally {
			loading = false;
		}
	}

	function onCountryChange() {
		regionCode = '';
		load();
	}
</script>

<section class="card">
	<h2>Where to look worldwide</h2>
	<p class="hint">
		Reference list of data sources EquiDrug uses per destination — local first, Google last (or
		skipped).
	</p>

	<div class="field">
		<label for="cat">Category</label>
		<select id="cat" bind:value={category} onchange={load}>
			<option value="drug">Drugs</option>
			<option value="supplement">Supplements</option>
			<option value="vitamin">Vitamins</option>
			<option value="food">Food / macros</option>
		</select>
	</div>

	<div class="field">
		<label for="region">Special region (optional)</label>
		<select id="region" bind:value={regionCode} onchange={load}>
			<option value="">— country below —</option>
			{#each specialRegions as r}
				<option value={r.code}>{r.name}</option>
			{/each}
		</select>
	</div>

	{#if !regionCode}
		<div class="field">
			<label for="country">Country</label>
			<select id="country" bind:value={countryCode} onchange={onCountryChange}>
				{#each indexed as c}
					<option value={c.code}>{c.name} ({c.google_tier})</option>
				{/each}
			</select>
		</div>
	{/if}

	<button class="btn btn-primary btn-block" onclick={load} disabled={loading}>Refresh</button>
</section>

{#if data}
	<DestinationSources countryCode={countryCode} regionCode={regionCode} {category} compact={false} />
{/if}

<section class="card">
	<h3>Indexed countries</h3>
	<ul class="index-list">
		{#each indexed as c}
			<li>
				<button
					type="button"
					class="index-btn"
					onclick={() => {
						countryCode = c.code;
						regionCode = '';
						load();
					}}
				>
					<strong>{c.name}</strong>
					<span class="tier">{c.google_tier}</span>
				</button>
			</li>
		{/each}
	</ul>
	<p class="hint">
		Full reference: <code>docs/DESTINATION_SOURCES.md</code> · Catalog:
		<code>backend/internal/sources/data/destination_sources.json</code>
	</p>
</section>

<style>
	.index-list {
		list-style: none;
		margin: 0;
		padding: 0;
	}

	.index-btn {
		display: flex;
		justify-content: space-between;
		width: 100%;
		padding: 0.5rem 0;
		border: none;
		border-top: 1px solid var(--color-border);
		background: none;
		cursor: pointer;
		text-align: left;
		font: inherit;
	}

	.index-btn:first-child {
		border-top: none;
	}

	.tier {
		font-size: 0.75rem;
		color: var(--color-muted);
		text-transform: uppercase;
	}

	code {
		font-size: 0.8rem;
	}
</style>

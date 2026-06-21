<script lang="ts">
	import { locale } from '$lib/stores/locale';
	import { t } from '$lib/i18n';
	import { lookup, countries, type LookupResponse } from '$lib/api/client';

	let query = $state('Allegra');
	let countryCode = $state('EE');
	let city = $state('');
	let loading = $state(false);
	let error = $state('');
	let result = $state<LookupResponse | null>(null);

	async function search() {
		loading = true;
		error = '';
		result = null;
		try {
			result = await lookup({
				query,
				country_code: countryCode,
				city: city || undefined,
				category: 'consume'
			});
		} catch (e) {
			error = e instanceof Error ? e.message : 'Search failed';
		} finally {
			loading = false;
		}
	}
</script>

<section class="card">
	<h2>{t($locale, 'lookup.title')}</h2>
	<p class="hint">{t($locale, 'lookup.hint')}</p>

	<div class="field">
		<label for="query">{t($locale, 'lookup.query')}</label>
		<input id="query" bind:value={query} placeholder={t($locale, 'lookup.example')} />
	</div>

	<div class="field">
		<label for="country">{t($locale, 'lookup.country')}</label>
		<select id="country" bind:value={countryCode}>
			{#each countries as c}
				<option value={c.code}>{c.name}</option>
			{/each}
		</select>
	</div>

	<div class="field">
		<label for="city">{t($locale, 'lookup.city')}</label>
		<input id="city" bind:value={city} placeholder="Tallinn" />
	</div>

	<button class="btn btn-primary btn-block" onclick={search} disabled={loading || !query.trim()}>
		{loading ? '…' : t($locale, 'lookup.search')}
	</button>

	{#if error}
		<p class="disclaimer" style="border-color: var(--color-danger); color: var(--color-danger);">
			{error}
		</p>
	{/if}
</section>

{#if result}
	<section class="card">
		<h2>Results for “{result.query}” in {result.country_code}</h2>
		{#if result.matches.length === 0}
			<div class="empty-state">No matches yet — curated data grows with each lookup.</div>
		{:else}
			{#each result.matches as match}
				<article class="match-card">
					<h3>{match.product.display_name}</h3>
					<p><span class="badge">{match.product.brand_name}</span> · {match.product.country_code}</p>
					{#if match.product.retailer_hint}
						<p>🛒 {match.product.retailer_hint}</p>
					{/if}
					{#if match.product.price_hint}
						<p>💰 {match.product.price_hint}</p>
					{/if}
					{#if match.equivalence.notes}
						<p>{match.equivalence.notes}</p>
					{/if}
					<p
						class={match.equivalence.confidence >= 0.9
							? 'confidence-high'
							: 'confidence-med'}
					>
						Confidence: {Math.round(match.equivalence.confidence * 100)}%
					</p>
				</article>
			{/each}
		{/if}
		<p class="disclaimer">{result.disclaimer}</p>
	</section>
{/if}

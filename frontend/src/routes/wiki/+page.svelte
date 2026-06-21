<script lang="ts">
	import { onMount } from 'svelte';
	import { locale } from '$lib/stores/locale';
	import { t } from '$lib/i18n';
	import { listWiki, type LookupMatch } from '$lib/api/client';

	let entries = $state<LookupMatch[]>([]);

	onMount(async () => {
		try {
			entries = await listWiki();
		} catch {
			entries = [];
		}
	});
</script>

<section class="card">
	<h2>{t($locale, 'wiki.title')}</h2>
	<p class="hint">{t($locale, 'wiki.empty')}</p>

	{#if entries.length === 0}
		<div class="empty-state">
			No saved equivalents yet. Run a lookup or convert a trip to grow your wiki.
		</div>
	{:else}
		{#each entries as entry}
			<article class="match-card">
				<h3>{entry.product.display_name}</h3>
				<p>{entry.product.country_code} · {Math.round(entry.equivalence.confidence * 100)}% match</p>
				{#if entry.equivalence.notes}
					<p>{entry.equivalence.notes}</p>
				{/if}
			</article>
		{/each}
	{/if}
</section>

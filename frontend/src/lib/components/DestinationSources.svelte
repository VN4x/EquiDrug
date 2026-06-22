<script lang="ts">
	import { onMount } from 'svelte';
	import {
		getDestinationSources,
		countries,
		type DestinationSourcesResponse
	} from '$lib/api/client';

	interface Props {
		countryCode: string;
		regionCode?: string;
		category?: string;
		compact?: boolean;
	}

	let { countryCode, regionCode = '', category = 'drug', compact = false }: Props = $props();

	let data = $state<DestinationSourcesResponse | null>(null);
	let loading = $state(false);
	let expanded = $state(true);

	$effect(() => {
		expanded = !compact;
	});

	$effect(() => {
		load(countryCode, regionCode, category);
	});

	async function load(cc: string, rc: string, cat: string) {
		if (!cc && !rc) return;
		loading = true;
		try {
			data = await getDestinationSources(cc, cat, rc || undefined);
		} catch {
			data = null;
		} finally {
			loading = false;
		}
	}

	function tierClass(tier: string) {
		if (tier === 'avoid') return 'tier-avoid';
		if (tier === 'limited') return 'tier-limited';
		if (tier === 'secondary') return 'tier-secondary';
		return 'tier-primary';
	}
</script>

{#if loading}
	<p class="sources-hint">Loading local sources…</p>
{:else if data}
	<section class="sources-panel" class:compact>
		<button type="button" class="sources-toggle" onclick={() => (expanded = !expanded)}>
			<span
				>Where we look in {data.region_name || data.country_name}</span
			>
			<span class="tier {tierClass(data.google_tier)}">{data.google_tier}</span>
			<span class="chevron">{expanded ? '▾' : '▸'}</span>
		</button>

		{#if expanded}
			<p class="guidance">{data.guidance}</p>
			<ol class="source-list">
				{#each data.sources as src}
					<li class:google={src.google}>
						<span class="prio">{src.priority}</span>
						<div class="src-body">
							<strong>
								{#if src.url.startsWith('http')}
									<a href={src.url} target="_blank" rel="noopener noreferrer">{src.name}</a>
								{:else}
									{src.name}
								{/if}
							</strong>
							<span class="type">{src.type_label ?? src.type}</span>
							{#if src.notes}
								<span class="note">{src.notes}</span>
							{/if}
						</div>
					</li>
				{/each}
			</ol>
			<p class="catalog-ver">Source catalog v{data.catalog_version}</p>
		{/if}
	</section>
{/if}

<style>
	.sources-panel {
		background: #f0f9f9;
		border: 1px solid var(--color-border);
		border-radius: 10px;
		padding: 0.75rem 1rem;
		margin-bottom: 1rem;
	}

	.sources-panel.compact {
		padding: 0.5rem 0.75rem;
	}

	.sources-toggle {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		width: 100%;
		background: none;
		border: none;
		cursor: pointer;
		font-weight: 600;
		font-size: 0.9rem;
		color: var(--color-primary);
		text-align: left;
		padding: 0;
	}

	.chevron {
		margin-left: auto;
		color: var(--color-muted);
	}

	.tier {
		font-size: 0.68rem;
		padding: 0.1rem 0.45rem;
		border-radius: 999px;
		font-weight: 600;
		text-transform: uppercase;
	}

	.tier-primary {
		background: #e8f8ee;
		color: var(--color-success);
	}
	.tier-secondary {
		background: #e8f2f1;
		color: var(--color-primary);
	}
	.tier-limited {
		background: #fef3e2;
		color: #b7791f;
	}
	.tier-avoid {
		background: #fde8e6;
		color: var(--color-danger);
	}

	.guidance {
		margin: 0.65rem 0 0.5rem;
		font-size: 0.82rem;
		color: var(--color-muted);
	}

	.source-list {
		margin: 0;
		padding: 0;
		list-style: none;
		max-height: 220px;
		overflow-y: auto;
	}

	.source-list li {
		display: flex;
		gap: 0.5rem;
		padding: 0.35rem 0;
		border-top: 1px solid rgba(0, 0, 0, 0.05);
		font-size: 0.82rem;
	}

	.source-list li.google {
		opacity: 0.75;
	}

	.prio {
		flex-shrink: 0;
		width: 1.25rem;
		color: var(--color-muted);
		font-size: 0.75rem;
	}

	.src-body {
		display: flex;
		flex-direction: column;
		gap: 0.1rem;
	}

	.type {
		color: var(--color-muted);
		font-size: 0.75rem;
	}

	.note {
		font-size: 0.72rem;
		color: #b7791f;
	}

	.catalog-ver {
		margin: 0.5rem 0 0;
		font-size: 0.68rem;
		color: var(--color-muted);
	}

	.sources-hint {
		font-size: 0.85rem;
		color: var(--color-muted);
	}
</style>

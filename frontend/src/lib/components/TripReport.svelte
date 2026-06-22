<script lang="ts">
	import {
		tripExportUrl,
		updateTripLineItem,
		countryName,
		formatDate,
		type TripReport
	} from '$lib/api/client';

	interface Props {
		report: TripReport;
		onupdate?: (report: TripReport) => void;
	}

	let { report, onupdate }: Props = $props();
	let saving = $state<string | null>(null);

	const tripId = $derived(report.trip.id);

	function productImage(url: string | undefined) {
		return url || '/products/allegra-us.svg';
	}

	async function toggleBought(rowId: string, bought: boolean) {
		saving = rowId;
		try {
			const updated = await updateTripLineItem(tripId, rowId, { bought });
			onupdate?.(updated);
		} finally {
			saving = null;
		}
	}

	async function saveNotes(rowId: string, notes: string) {
		saving = rowId;
		try {
			const updated = await updateTripLineItem(tripId, rowId, { notes });
			onupdate?.(updated);
		} finally {
			saving = null;
		}
	}

	function openExport() {
		window.open(tripExportUrl(tripId, 'html'), '_blank');
	}

	function printReport() {
		window.open(tripExportUrl(tripId, 'html'), '_blank')?.print();
	}
</script>

<article class="trip-report">
	<header class="report-hero">
		<div>
			<h2>{report.trip.title}</h2>
			<p class="route">
				{countryName(report.trip.origin_country)} → {countryName(report.trip.dest_country)}
				{#if report.trip.dest_city}
					· {report.trip.dest_city}
				{/if}
			</p>
			<p class="dates">
				{formatDate(report.trip.start_date)} – {formatDate(report.trip.end_date)}
				· {report.summary.trip_days} days · +{report.trip.spare_percent}% spare
			</p>
		</div>
		<div class="export-actions">
			<button class="btn btn-secondary" type="button" onclick={openExport}>Export HTML</button>
			<button class="btn btn-secondary" type="button" onclick={printReport}>Print / PDF</button>
		</div>
	</header>

	<div class="stats-grid">
		<div class="stat">
			<strong>{report.summary.total_items}</strong>
			<span>Items</span>
		</div>
		<div class="stat">
			<strong>{report.summary.matched_items}</strong>
			<span>Matched</span>
		</div>
		<div class="stat">
			<strong>{report.summary.unmatched_items}</strong>
			<span>To find</span>
		</div>
		<div class="stat">
			<strong>{report.summary.bought_items}</strong>
			<span>Bought</span>
		</div>
	</div>

	<section class="equivalence-list">
		<h3>Equivalence report</h3>
		{#each report.rows as row (row.line_item.id)}
			<div class="equiv-row" class:bought={row.line_item.bought}>
				<div class="equiv-header">
					<span class="badge">{row.locker_item.custom_name}</span>
					<span class="qty"
						>Need <strong>{row.line_item.quantity_needed} {row.line_item.quantity_unit}</strong></span
					>
					{#if row.match_status === 'matched'}
						<span class="status matched">{Math.round(row.confidence * 100)}% match</span>
					{:else}
						<span class="status unmatched">No match yet</span>
					{/if}
				</div>

				<div class="product-pair">
					<div class="product-card home">
						<p class="label">At home</p>
						{#if row.origin}
							<img src={productImage(row.origin.image_url)} alt={row.origin.display_name} />
							<h4>{row.origin.display_name}</h4>
							<p class="meta">{row.origin.brand_name} · {row.origin.country_code}</p>
						{:else}
							<div class="placeholder-img">?</div>
							<h4>{row.locker_item.custom_name}</h4>
						{/if}
					</div>

					<div class="arrow" aria-hidden="true">→</div>

					<div class="product-card foreign" class:empty={!row.foreign}>
						<p class="label">Buy locally</p>
						{#if row.foreign}
							<img src={productImage(row.foreign.image_url)} alt={row.foreign.display_name} />
							<h4>{row.foreign.display_name}</h4>
							<p class="meta">{row.foreign.brand_name} · {row.foreign.country_code}</p>
						{:else}
							<div class="placeholder-img warn">!</div>
							<h4>Lookup at destination</h4>
						{/if}
					</div>
				</div>

				{#if row.foreign}
					<div class="buy-hints">
						{#if row.line_item.where_to_buy}
							<p>🛒 {row.line_item.where_to_buy}</p>
						{/if}
						{#if row.line_item.estimated_price}
							<p>💰 {row.line_item.estimated_price}</p>
						{/if}
						{#if row.match_notes}
							<p class="notes">{row.match_notes}</p>
						{/if}
					</div>
				{:else if row.match_notes}
					<p class="notes">{row.match_notes}</p>
				{/if}
			</div>
		{/each}
	</section>

	<section class="shopping-list card">
		<h3>Shopping list</h3>
		<p class="hint">Mark items when bought and add notes (pharmacy found, lot number, etc.)</p>

		{#each report.rows as row (row.line_item.id)}
			<div class="shop-item" class:bought={row.line_item.bought}>
				<label class="check-row">
					<input
						type="checkbox"
						checked={row.line_item.bought}
						disabled={saving === row.line_item.id}
						onchange={(e) =>
							toggleBought(row.line_item.id, (e.currentTarget as HTMLInputElement).checked)}
					/>
					<span class="shop-label">
						<strong>{row.foreign?.display_name ?? row.locker_item.custom_name}</strong>
						<span class="shop-qty"
							>{row.line_item.quantity_needed} {row.line_item.quantity_unit}</span
						>
					</span>
				</label>
				<textarea
					class="shop-notes"
					placeholder="Notes: pharmacy name, aisle, paid amount…"
					value={row.line_item.notes ?? ''}
					disabled={saving === row.line_item.id}
					onblur={(e) => {
						const v = (e.currentTarget as HTMLTextAreaElement).value;
						if (v !== (row.line_item.notes ?? '')) saveNotes(row.line_item.id, v);
					}}
				></textarea>
			</div>
		{/each}
	</section>

	<p class="disclaimer">{report.disclaimer}</p>
</article>

<style>
	.trip-report {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.report-hero {
		background: linear-gradient(135deg, #0d6e6e, #095555);
		color: #fff;
		border-radius: var(--radius);
		padding: 1.25rem;
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.report-hero h2 {
		margin: 0;
		font-size: 1.35rem;
	}

	.route,
	.dates {
		margin: 0.25rem 0 0;
		opacity: 0.92;
		font-size: 0.9rem;
	}

	.export-actions {
		display: flex;
		gap: 0.5rem;
		flex-wrap: wrap;
	}

	.export-actions .btn {
		background: rgba(255, 255, 255, 0.15);
		color: #fff;
		border: 1px solid rgba(255, 255, 255, 0.35);
	}

	.stats-grid {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 0.5rem;
	}

	.stat {
		background: var(--color-surface);
		border-radius: 10px;
		padding: 0.85rem 0.5rem;
		text-align: center;
		border: 1px solid var(--color-border);
	}

	.stat strong {
		display: block;
		font-size: 1.4rem;
		color: var(--color-primary);
	}

	.stat span {
		font-size: 0.75rem;
		color: var(--color-muted);
	}

	.equivalence-list h3,
	.shopping-list h3 {
		margin: 0 0 0.75rem;
	}

	.equiv-row {
		background: var(--color-surface);
		border: 1px solid var(--color-border);
		border-radius: var(--radius);
		padding: 1rem;
		margin-bottom: 0.75rem;
	}

	.equiv-row.bought {
		opacity: 0.72;
	}

	.equiv-header {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
		align-items: center;
		margin-bottom: 0.85rem;
	}

	.qty {
		font-size: 0.9rem;
		color: var(--color-muted);
	}

	.status {
		margin-left: auto;
		font-size: 0.8rem;
		font-weight: 600;
		padding: 0.15rem 0.5rem;
		border-radius: 999px;
	}

	.status.matched {
		background: #e8f8ee;
		color: var(--color-success);
	}

	.status.unmatched {
		background: #fef3e2;
		color: #b7791f;
	}

	.product-pair {
		display: grid;
		grid-template-columns: 1fr auto 1fr;
		gap: 0.75rem;
		align-items: center;
	}

	.product-card {
		text-align: center;
		border: 1px solid var(--color-border);
		border-radius: 10px;
		padding: 0.75rem;
		background: #fafcfc;
	}

	.product-card img {
		width: 72px;
		height: 72px;
		object-fit: contain;
		margin: 0 auto 0.5rem;
		display: block;
	}

	.product-card h4 {
		margin: 0 0 0.25rem;
		font-size: 0.85rem;
		line-height: 1.3;
	}

	.product-card .label {
		margin: 0 0 0.5rem;
		font-size: 0.72rem;
		text-transform: uppercase;
		letter-spacing: 0.04em;
		color: var(--color-muted);
	}

	.product-card .meta {
		margin: 0;
		font-size: 0.75rem;
		color: var(--color-muted);
	}

	.placeholder-img {
		width: 72px;
		height: 72px;
		margin: 0 auto 0.5rem;
		border-radius: 8px;
		background: #e8f2f1;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1.5rem;
		color: var(--color-primary);
	}

	.placeholder-img.warn {
		background: #fef3e2;
		color: #b7791f;
	}

	.arrow {
		font-size: 1.5rem;
		color: var(--color-primary);
		text-align: center;
	}

	.buy-hints p,
	.notes {
		margin: 0.35rem 0 0;
		font-size: 0.85rem;
		color: var(--color-muted);
	}

	.shop-item {
		border-top: 1px solid var(--color-border);
		padding: 0.85rem 0;
	}

	.shop-item:first-of-type {
		border-top: none;
	}

	.shop-item.bought .shop-label {
		text-decoration: line-through;
		opacity: 0.7;
	}

	.check-row {
		display: flex;
		align-items: flex-start;
		gap: 0.65rem;
		cursor: pointer;
	}

	.check-row input {
		margin-top: 0.2rem;
		width: 1.1rem;
		height: 1.1rem;
		accent-color: var(--color-primary);
	}

	.shop-label {
		display: flex;
		flex-direction: column;
		gap: 0.15rem;
	}

	.shop-qty {
		font-size: 0.85rem;
		color: var(--color-muted);
	}

	.shop-notes {
		width: 100%;
		margin-top: 0.5rem;
		padding: 0.5rem 0.65rem;
		border: 1px solid var(--color-border);
		border-radius: 8px;
		font-size: 0.85rem;
		font-family: inherit;
		min-height: 2.5rem;
		resize: vertical;
	}

	@media (max-width: 560px) {
		.stats-grid {
			grid-template-columns: repeat(2, 1fr);
		}

		.product-pair {
			grid-template-columns: 1fr;
		}

		.arrow {
			transform: rotate(90deg);
		}

		.status {
			margin-left: 0;
		}
	}
</style>

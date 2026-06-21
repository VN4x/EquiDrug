<script lang="ts">
	interface Props {
		label: string;
		consumed: number;
		target: number;
		unit?: string;
		color: string;
	}

	let { label, consumed, target, unit = 'g', color }: Props = $props();

	const pct = $derived(target > 0 ? Math.min(100, (consumed / target) * 100) : 0);
	const remaining = $derived(Math.max(0, target - consumed));
	const over = $derived(consumed > target);
</script>

<div class="macro-bar">
	<div class="macro-bar-head">
		<span class="label">{label}</span>
		<span class="nums" class:over>
			<strong>{consumed.toFixed(0)}</strong> / {target}{unit}
		</span>
	</div>
	<div class="track">
		<div
			class="fill"
			style="width: {pct}%; background: {over ? 'var(--color-danger)' : color}"
		></div>
	</div>
	<p class="remain">
		{#if over}
			{Math.abs(target - consumed).toFixed(0)}{unit} over target
		{:else}
			{remaining.toFixed(0)}{unit} left
		{/if}
	</p>
</div>

<style>
	.macro-bar {
		margin-bottom: 0.85rem;
	}

	.macro-bar-head {
		display: flex;
		justify-content: space-between;
		align-items: baseline;
		margin-bottom: 0.35rem;
		font-size: 0.85rem;
	}

	.label {
		font-weight: 600;
	}

	.nums {
		color: var(--color-muted);
	}

	.nums.over {
		color: var(--color-danger);
	}

	.track {
		height: 10px;
		background: #e8f2f1;
		border-radius: 999px;
		overflow: hidden;
	}

	.fill {
		height: 100%;
		border-radius: 999px;
		transition: width 0.3s ease;
	}

	.remain {
		margin: 0.25rem 0 0;
		font-size: 0.75rem;
		color: var(--color-muted);
	}
</style>

<script lang="ts">
	import { onMount } from 'svelte';
	import { locale } from '$lib/stores/locale';
	import { t } from '$lib/i18n';
	import { addLockerItem, listLocker, scanLocker, type LockerItem } from '$lib/api/client';

	let items = $state<LockerItem[]>([]);
	let customName = $state('');
	let dosePerDay = $state(1);
	let loading = $state(false);
	let message = $state('');

	onMount(async () => {
		try {
			items = await listLocker();
		} catch {
			items = [];
		}
	});

	async function addItem() {
		if (!customName.trim()) return;
		loading = true;
		try {
			const created = await addLockerItem({
				custom_name: customName,
				dose_per_day: dosePerDay,
				dose_unit: 'tablet',
				frequency: 'daily',
				category: 'consume'
			});
			items = [created, ...items];
			customName = '';
			message = 'Added to locker';
		} catch (e) {
			message = e instanceof Error ? e.message : 'Failed';
		} finally {
			loading = false;
		}
	}

	async function scan() {
		loading = true;
		message = '';
		try {
			const res = await scanLocker('placeholder://locker-photo');
			message = res.message;
		} catch (e) {
			message = e instanceof Error ? e.message : 'Scan failed';
		} finally {
			loading = false;
		}
	}
</script>

<section class="card">
	<h2>{t($locale, 'locker.title')}</h2>
	<p class="hint">{t($locale, 'locker.hint')}</p>

	<button class="btn btn-secondary btn-block" onclick={scan} disabled={loading}>
		📷 {t($locale, 'locker.scan')}
	</button>

	<div class="field" style="margin-top: 1rem;">
		<label for="name">Item name</label>
		<input id="name" bind:value={customName} placeholder="Allegra 180mg" />
	</div>
	<div class="field">
		<label for="dose">Dose per day</label>
		<input id="dose" type="number" min="0.5" step="0.5" bind:value={dosePerDay} />
	</div>
	<button class="btn btn-primary btn-block" onclick={addItem} disabled={loading || !customName.trim()}>
		{t($locale, 'locker.add')}
	</button>

	{#if message}
		<p class="disclaimer">{message}</p>
	{/if}
</section>

<section class="card">
	<h2>Inventory</h2>
	{#if items.length === 0}
		<div class="empty-state">Your locker is empty. Add items or scan a photo.</div>
	{:else}
		{#each items as item}
			<article class="match-card">
				<h3>{item.custom_name}</h3>
				<p>{item.dose_per_day} {item.dose_unit} · {item.frequency}</p>
			</article>
		{/each}
	{/if}
</section>

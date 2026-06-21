<script lang="ts">
	import { locale } from '$lib/stores/locale';
	import { t } from '$lib/i18n';
	import TripReport from '$lib/components/TripReport.svelte';
	import { convertTrip, createTrip, countries, type TripReport as TripReportType } from '$lib/api/client';

	let title = $state('Baltic spring trip');
	let originCountry = $state('US');
	let destCountry = $state('EE');
	let destCity = $state('Tallinn');
	let startDate = $state('2026-04-01');
	let endDate = $state('2026-04-14');
	let sparePercent = $state(5);
	let preferredBrands = $state('');
	let loading = $state(false);
	let report = $state<TripReportType | null>(null);
	let error = $state('');

	async function plan() {
		loading = true;
		error = '';
		report = null;
		try {
			const trip = await createTrip({
				title,
				origin_country: originCountry,
				dest_country: destCountry,
				dest_city: destCity,
				start_date: startDate,
				end_date: endDate,
				spare_percent: sparePercent,
				preferred_brands: preferredBrands
					.split(',')
					.map((s) => s.trim())
					.filter(Boolean)
			});
			const brands = preferredBrands
				.split(',')
				.map((s) => s.trim())
				.filter(Boolean);
			report = await convertTrip(trip.id, brands.length ? brands : undefined);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Planning failed';
		} finally {
			loading = false;
		}
	}

	function onReportUpdate(updated: TripReportType) {
		report = updated;
	}
</script>

<section class="card">
	<h2>{t($locale, 'planner.title')}</h2>
	<p class="hint">{t($locale, 'planner.hint')}</p>
	<p class="hint" style="margin-top: -0.5rem;">
		Demo locker has Allegra + Zyrtec — try US → EE or US → DE.
	</p>

	<div class="field">
		<label for="title">Trip name</label>
		<input id="title" bind:value={title} />
	</div>

	<div class="field">
		<label for="origin">From (home country)</label>
		<select id="origin" bind:value={originCountry}>
			{#each countries as c}
				<option value={c.code}>{c.name}</option>
			{/each}
		</select>
	</div>

	<div class="field">
		<label for="dest">Destination</label>
		<select id="dest" bind:value={destCountry}>
			{#each countries as c}
				<option value={c.code}>{c.name}</option>
			{/each}
		</select>
	</div>

	<div class="field">
		<label for="city">City</label>
		<input id="city" bind:value={destCity} />
	</div>

	<div class="field">
		<label for="start">Start date</label>
		<input id="start" type="date" bind:value={startDate} />
	</div>

	<div class="field">
		<label for="end">End date</label>
		<input id="end" type="date" bind:value={endDate} />
	</div>

	<div class="field">
		<label for="spare">Spare stock %</label>
		<input id="spare" type="number" min="0" max="50" bind:value={sparePercent} />
	</div>

	<div class="field">
		<label for="brands">Preferred retailers (comma-separated)</label>
		<input id="brands" bind:value={preferredBrands} placeholder="Benu, Apotheka" />
	</div>

	<button class="btn btn-primary btn-block" onclick={plan} disabled={loading}>
		{loading ? 'Converting…' : t($locale, 'planner.convert')}
	</button>

	{#if error}
		<p class="disclaimer" style="border-color: var(--color-danger); color: var(--color-danger);">
			{error}
		</p>
	{/if}
</section>

{#if report}
	<TripReport {report} onupdate={onReportUpdate} />
{/if}

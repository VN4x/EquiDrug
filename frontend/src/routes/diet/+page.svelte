<script lang="ts">
	import { onMount } from 'svelte';
	import { locale } from '$lib/stores/locale';
	import { t } from '$lib/i18n';
	import MacroBar from '$lib/components/MacroBar.svelte';
	import {
		analyzeFood,
		deleteFoodLog,
		getDietToday,
		logFood,
		updateDietProfile,
		countries,
		type AnalyzeFoodResponse,
		type DietDaySummary,
		type FoodReference
	} from '$lib/api/client';

	let day = $state<DietDaySummary | null>(null);
	let showTargets = $state(false);
	let proteinTarget = $state(130);
	let carbsTarget = $state(80);
	let fatTarget = $state(25);

	let foodQuery = $state('');
	let menuText = $state('');
	let countryCode = $state('VN');
	let photoPreview = $state<string | null>(null);
	let hasPhoto = $state(false);

	let analysis = $state<AnalyzeFoodResponse | null>(null);
	let selected = $state<FoodReference | null>(null);
	let loading = $state(false);
	let message = $state('');

	onMount(() => loadDay());

	async function loadDay() {
		try {
			day = await getDietToday();
			proteinTarget = day.profile.protein_g;
			carbsTarget = day.profile.carbs_g;
			fatTarget = day.profile.fat_g;
		} catch {
			day = null;
		}
	}

	async function saveTargets() {
		loading = true;
		try {
			await updateDietProfile({
				protein_g: proteinTarget,
				carbs_g: carbsTarget,
				fat_g: fatTarget
			});
			await loadDay();
			showTargets = false;
			message = 'Targets saved';
		} catch (e) {
			message = e instanceof Error ? e.message : 'Save failed';
		} finally {
			loading = false;
		}
	}

	function onPhotoSelect(e: Event) {
		const input = e.currentTarget as HTMLInputElement;
		const file = input.files?.[0];
		if (!file) return;
		hasPhoto = true;
		const reader = new FileReader();
		reader.onload = () => {
			photoPreview = reader.result as string;
		};
		reader.readAsDataURL(file);
	}

	async function runAnalyze() {
		loading = true;
		message = '';
		analysis = null;
		selected = null;
		try {
			analysis = await analyzeFood({
				query: foodQuery || undefined,
				menu_text: menuText || undefined,
				image_url: hasPhoto ? 'local://plate-photo' : undefined,
				country_code: countryCode
			});
			if (analysis.suggested) {
				selected = analysis.suggested;
			}
		} catch (e) {
			message = e instanceof Error ? e.message : 'Analyze failed';
		} finally {
			loading = false;
		}
	}

	function pickMatch(ref: FoodReference) {
		selected = ref;
	}

	async function addSelected() {
		if (!selected) return;
		loading = true;
		try {
			day = await logFood({
				food_name: selected.name,
				serving_label: selected.serving_label,
				protein_g: selected.protein_g,
				carbs_g: selected.carbs_g,
				fat_g: selected.fat_g,
				calories_kcal: selected.calories_kcal,
				confidence: 0.9,
				source: 'analyze',
				notes: selected.name_local ? `Local: ${selected.name_local}` : undefined
			});
			analysis = null;
			selected = null;
			foodQuery = '';
			menuText = '';
			hasPhoto = false;
			photoPreview = null;
			message = 'Added to today';
		} catch (e) {
			message = e instanceof Error ? e.message : 'Log failed';
		} finally {
			loading = false;
		}
	}

	async function removeEntry(id: string) {
		if (!day) return;
		loading = true;
		try {
			day = await deleteFoodLog(id, day.date);
		} finally {
			loading = false;
		}
	}
</script>

<section class="card diet-hero">
	<h2>{t($locale, 'diet.title')}</h2>
	<p class="hint">{t($locale, 'diet.hint')}</p>
	<p class="targets-line">
		{t($locale, 'diet.targets')}: <strong>{proteinTarget}g P</strong> ·
		<strong>{carbsTarget}g C</strong> · <strong>{fatTarget}g F</strong>
		<button class="link-btn" type="button" onclick={() => (showTargets = !showTargets)}>Edit</button>
	</p>
</section>

{#if showTargets}
	<section class="card">
		<h3>{t($locale, 'diet.targets')}</h3>
		<div class="field">
			<label for="p">{t($locale, 'diet.protein')} (g)</label>
			<input id="p" type="number" min="0" bind:value={proteinTarget} />
		</div>
		<div class="field">
			<label for="c">{t($locale, 'diet.carbs')} (g)</label>
			<input id="c" type="number" min="0" bind:value={carbsTarget} />
		</div>
		<div class="field">
			<label for="f">{t($locale, 'diet.fat')} (g)</label>
			<input id="f" type="number" min="0" bind:value={fatTarget} />
		</div>
		<button class="btn btn-primary btn-block" onclick={saveTargets} disabled={loading}>Save targets</button>
	</section>
{/if}

{#if day}
	<section class="card macro-dashboard">
		<h3>Today · {day.date}</h3>
		<MacroBar
			label={t($locale, 'diet.protein')}
			consumed={day.consumed.protein_g}
			target={day.profile.protein_g}
			color="#0d6e6e"
		/>
		<MacroBar
			label={t($locale, 'diet.carbs')}
			consumed={day.consumed.carbs_g}
			target={day.profile.carbs_g}
			color="#e8a838"
		/>
		<MacroBar
			label={t($locale, 'diet.fat')}
			consumed={day.consumed.fat_g}
			target={day.profile.fat_g}
			color="#5c7cfa"
		/>
		<p class="kcal-meta">{day.consumed.calories_kcal.toFixed(0)} kcal consumed</p>
	</section>
{/if}

<section class="card">
	<h3>Log food</h3>

	<label class="photo-btn">
		<input type="file" accept="image/*" capture="environment" hidden onchange={onPhotoSelect} />
		📷 {t($locale, 'diet.photo')}
	</label>
	{#if photoPreview}
		<img class="plate-preview" src={photoPreview} alt="Plate preview" />
	{/if}

	<div class="field">
		<label for="food">Dish name</label>
		<input
			id="food"
			bind:value={foodQuery}
			placeholder="pho bo, pad thai, grilled chicken…"
		/>
	</div>

	<div class="field">
		<label for="menu">{t($locale, 'diet.menu')}</label>
		<textarea
			id="menu"
			bind:value={menuText}
			rows="2"
			placeholder="Paste menu line: Phở bò tái — beef noodle soup"
		></textarea>
	</div>

	<div class="field">
		<label for="cc">Country (prioritize local dishes)</label>
		<select id="cc" bind:value={countryCode}>
			{#each countries as c}
				<option value={c.code}>{c.name}</option>
			{/each}
		</select>
	</div>

	<button
		class="btn btn-primary btn-block"
		onclick={runAnalyze}
		disabled={loading || (!foodQuery.trim() && !menuText.trim() && !hasPhoto)}
	>
		{loading ? '…' : t($locale, 'diet.analyze')}
	</button>
</section>

{#if analysis}
	<section class="card analysis-card">
		{#if analysis.vision_note}
			<p class="vision-note">{analysis.vision_note}</p>
		{/if}

		{#if analysis.matches.length === 0}
			<p class="empty-hint">No match — enter a clearer name or log macros manually later.</p>
		{:else}
			<h3>Matches</h3>
			{#each analysis.matches as match}
				<button
					type="button"
					class="food-match"
					class:selected={selected?.id === match.reference.id}
					onclick={() => pickMatch(match.reference)}
				>
					<div class="match-top">
						<strong>{match.reference.name}</strong>
						<span class="badge">{Math.round(match.confidence * 100)}%</span>
					</div>
					{#if match.reference.name_local}
						<p class="local">{match.reference.name_local}</p>
					{/if}
					<p class="macros">
						P {match.reference.protein_g}g · C {match.reference.carbs_g}g · F
						{match.reference.fat_g}g · {match.reference.serving_label}
					</p>
					<p class="reason">{match.match_reason}</p>
				</button>
			{/each}

			{#if selected}
				<button class="btn btn-primary btn-block" onclick={addSelected} disabled={loading}>
					{t($locale, 'diet.add')}: {selected.name}
				</button>
			{/if}
		{/if}

		<p class="disclaimer">{analysis.disclaimer}</p>
	</section>
{/if}

{#if day && day.entries.length > 0}
	<section class="card">
		<h3>{t($locale, 'diet.today')}</h3>
		{#each day.entries as entry (entry.id)}
			<article class="log-entry">
				<div class="log-head">
					<strong>{entry.food_name}</strong>
					<button
						type="button"
						class="link-btn danger"
						onclick={() => removeEntry(entry.id)}
						disabled={loading}>Remove</button
					>
				</div>
				<p class="macros">
					P {entry.protein_g}g · C {entry.carbs_g}g · F {entry.fat_g}g
					{#if entry.serving_label}
						· {entry.serving_label}
					{/if}
				</p>
				{#if entry.notes}
					<p class="notes">{entry.notes}</p>
				{/if}
			</article>
		{/each}
	</section>
{/if}

{#if message}
	<p class="disclaimer">{message}</p>
{/if}

<style>
	.diet-hero h2 {
		margin-bottom: 0.25rem;
	}

	.targets-line {
		margin: 0.75rem 0 0;
		font-size: 0.9rem;
		color: var(--color-muted);
	}

	.link-btn {
		background: none;
		border: none;
		color: var(--color-primary);
		font-weight: 600;
		cursor: pointer;
		font-size: inherit;
		padding: 0 0.25rem;
	}

	.link-btn.danger {
		color: var(--color-danger);
		font-size: 0.8rem;
	}

	.macro-dashboard h3 {
		margin: 0 0 0.75rem;
	}

	.kcal-meta {
		margin: 0;
		font-size: 0.85rem;
		color: var(--color-muted);
		text-align: right;
	}

	.photo-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		padding: 0.85rem;
		border: 2px dashed var(--color-border);
		border-radius: 10px;
		cursor: pointer;
		margin-bottom: 1rem;
		font-weight: 600;
		color: var(--color-primary);
		background: #fafcfc;
	}

	.plate-preview {
		width: 100%;
		max-height: 180px;
		object-fit: cover;
		border-radius: 10px;
		margin-bottom: 1rem;
	}

	textarea {
		width: 100%;
		padding: 0.65rem 0.75rem;
		border: 1px solid var(--color-border);
		border-radius: 8px;
		font-size: 1rem;
		font-family: inherit;
		resize: vertical;
	}

	.vision-note {
		background: #fef9e8;
		border-radius: 8px;
		padding: 0.65rem 0.75rem;
		font-size: 0.85rem;
		margin: 0 0 0.75rem;
	}

	.food-match {
		display: block;
		width: 100%;
		text-align: left;
		border: 2px solid var(--color-border);
		border-radius: 10px;
		padding: 0.85rem;
		margin-bottom: 0.5rem;
		background: #fff;
		cursor: pointer;
	}

	.food-match.selected {
		border-color: var(--color-primary);
		background: #f0f9f9;
	}

	.match-top {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 0.5rem;
	}

	.local {
		margin: 0.25rem 0;
		font-size: 0.85rem;
		color: var(--color-primary);
	}

	.macros,
	.reason,
	.notes {
		margin: 0.2rem 0 0;
		font-size: 0.82rem;
		color: var(--color-muted);
	}

	.log-entry {
		border-top: 1px solid var(--color-border);
		padding: 0.75rem 0;
	}

	.log-entry:first-of-type {
		border-top: none;
	}

	.log-head {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 0.5rem;
	}

	.empty-hint {
		color: var(--color-muted);
		font-size: 0.9rem;
	}
</style>

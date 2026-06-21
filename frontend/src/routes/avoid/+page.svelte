<script lang="ts">
	import { locale } from '$lib/stores/locale';
	import { t } from '$lib/i18n';

	const sampleTraps = [
		{
			from: 'en / peanut',
			to: 'th / ถั่วลิสง',
			severity: 'critical',
			note: 'Ask for ปลอดถั่วลิสง when dining in Thailand'
		},
		{
			from: 'en / corn syrup',
			to: 'de / Glukosesirup',
			severity: 'info',
			note: 'Not always maize-derived — check label'
		},
		{
			from: 'en / gluten-free',
			to: 'fr / sans gluten',
			severity: 'info',
			note: 'EU-regulated claim; still verify certification'
		},
		{
			from: 'ru / арахис',
			to: 'en / peanut',
			severity: 'critical',
			note: 'Cyrillic ↔ Latin labeling trap'
		}
	];
</script>

<section class="card">
	<h2>{t($locale, 'avoid.title')}</h2>
	<p class="hint">{t($locale, 'avoid.hint')}</p>

	<p>
		Build your avoid list from diet photos (OCR+AI) or manual entry — same flow as the medicine
		locker, but tagged <span class="badge">{t($locale, 'category.avoid')}</span>.
	</p>

	{#each sampleTraps as trap}
		<article class="match-card">
			<h3>{trap.from} → {trap.to}</h3>
			<p>
				<span
					class="badge"
					style={trap.severity === 'critical'
						? 'background:#fde8e6;color:var(--color-danger)'
						: ''}
				>
					{trap.severity}
				</span>
			</p>
			<p>{trap.note}</p>
		</article>
	{/each}

	<p class="disclaimer">{t($locale, 'disclaimer')}</p>
</section>

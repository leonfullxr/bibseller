<script lang="ts">
	import { getI18n } from '$lib/i18n';

	// The price/face-value/description inputs shared by the create and edit
	// listing forms. The parent owns the <form>, submit button, and feedback.
	let {
		price = '',
		original = '',
		description = ''
	}: { price?: string; original?: string; description?: string } = $props();

	const { t } = getI18n();
</script>

<label for="price">{t('listingFields.price')}</label>
<!-- max mirrors decision D2 (ask <= face value) for pre-submit UX only; the
     server stays authoritative. -->
<input
	id="price"
	name="price"
	type="number"
	min="0"
	max={original || undefined}
	step="0.01"
	inputmode="decimal"
	class="field"
	value={price}
	placeholder={t('listingFields.pricePlaceholder')}
/>

<label for="original_price">{t('listingFields.original')}</label>
<input
	id="original_price"
	name="original_price"
	type="number"
	min="0"
	step="0.01"
	inputmode="decimal"
	class="field"
	value={original}
	placeholder={t('listingFields.optional')}
/>
<p class="hint">{t('listingFields.hint')}</p>

<label for="description">{t('listingFields.description')}</label>
<textarea
	id="description"
	name="description"
	rows="3"
	maxlength="2000"
	class="field"
	placeholder={t('listingFields.descriptionPlaceholder')}>{description}</textarea
>

<style>
	label {
		margin-top: 0.5rem;
		font-size: 0.75rem;
		line-height: 1rem;
		font-weight: 500;
		color: var(--slate-600);
	}

	.field {
		width: 100%;
	}

	textarea {
		resize: vertical;
	}

	.hint {
		font-size: 0.75rem;
		line-height: 1rem;
		color: var(--slate-500);
	}
</style>

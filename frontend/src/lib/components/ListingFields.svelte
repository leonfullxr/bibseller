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
<input
	id="price"
	name="price"
	type="number"
	min="0"
	step="0.01"
	inputmode="decimal"
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

	input,
	textarea {
		width: 100%;
		border-radius: 0.375rem;
		border: 1px solid var(--slate-300);
		background: white;
		padding: 0.5rem 0.75rem;
		font-size: 0.875rem;
		line-height: 1.25rem;
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

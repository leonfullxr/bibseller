<script lang="ts">
	// The price/face-value/description inputs shared by the create and edit
	// listing forms. The parent owns the <form>, submit button, and feedback.
	let {
		price = '',
		original = '',
		description = ''
	}: { price?: string; original?: string; description?: string } = $props();
</script>

<label for="price">Asking price (EUR)</label>
<input
	id="price"
	name="price"
	type="number"
	min="0"
	step="0.01"
	inputmode="decimal"
	value={price}
	placeholder="e.g. 45"
/>

<label for="original_price">Original price / face value (EUR)</label>
<input
	id="original_price"
	name="original_price"
	type="number"
	min="0"
	step="0.01"
	inputmode="decimal"
	value={original}
	placeholder="optional"
/>
<p class="hint">Enter the face value and your asking price is capped at it - no scalping.</p>

<label for="description">Description</label>
<textarea
	id="description"
	name="description"
	rows="3"
	maxlength="2000"
	placeholder="optional - size, pickup details, etc.">{description}</textarea
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

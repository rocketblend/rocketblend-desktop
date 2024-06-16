<script lang="ts">
	import type { SvelteComponent } from 'svelte';

	import { getModalStore } from '@skeletonlabs/skeleton';
	import { Media2 } from '$lib/components/ui/gallery';

	export let parent: SvelteComponent;

	const modalStore = getModalStore();

	const cButton = 'fixed top-4 right-4 z-50 font-bold shadow-xl';
	const cImage = 'max-w-[90%] max-h-[90%] rounded-container-token overflow-hidden shadow-xl';
</script>

{#if $modalStore[0]}
	<button class="btn-icon variant-filled {cButton}" on:click={parent.onClose}>×</button>
	{#if $modalStore[0]?.meta.src}
		<Media2 src={$modalStore[0]?.meta.src} class={cImage} alt="{$modalStore[0]?.meta.alt}"/>
	{/if}
{/if}
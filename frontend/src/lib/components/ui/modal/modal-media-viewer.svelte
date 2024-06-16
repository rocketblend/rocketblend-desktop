<script lang="ts">
    import type { SvelteComponent } from 'svelte';
    import { getModalStore } from '@skeletonlabs/skeleton';

	import IconClose from '~icons/ri/close-fill';
	import IconArrowRight from '~icons/ri/arrow-drop-right-line';
	import IconArrowLeft from '~icons/ri/arrow-drop-left-line';

    import { Media2 } from '$lib/components/ui/gallery';

    export let parent: SvelteComponent;

    const modalStore = getModalStore();

    const cButton = 'fixed top-4 right-4 z-50 shadow-xl';
    const cImage = 'overflow-hidden shadow-xl';

    let elemMedia: HTMLDivElement;

    function scrollLeft(): void {
        let x = elemMedia.scrollWidth;
        if (elemMedia.scrollLeft !== 0) x = elemMedia.scrollLeft - elemMedia.clientWidth;
        elemMedia.scroll({ left: x, behavior: 'smooth' });
    }

    function scrollRight(): void {
        let x = 0;
        if (elemMedia.scrollLeft < elemMedia.scrollWidth - elemMedia.clientWidth - 1) 
            x = elemMedia.scrollLeft + elemMedia.clientWidth;
        elemMedia.scroll({ left: x, behavior: 'smooth' });
    }
</script>


{#if $modalStore[0]}
<div class="modal block h-96 w-auto p-4 space-y-4">
    <button class="btn-icon variant-filled {cButton}" on:click={parent.onClose}><IconClose/></button>
    <div class="grid grid-cols-[auto_1fr_auto] gap-4 items-center">
        <button type="button" class="btn-icon variant-filled text-xl" on:click={scrollLeft}>
            <IconArrowLeft />
        </button>
        <div bind:this={elemMedia} class="snap-x snap-mandatory scroll-smooth flex gap-2 pb-2 overflow-x-auto">
            {#each $modalStore[0]?.meta.media as mediaItem}
                <div class="shrink-0 h-full w-full snap-start">
                    <Media2 src={mediaItem.url} class={cImage} rounded />
                </div>
            {/each}
        </div>
        <button type="button" class="btn-icon variant-filled text-xl" on:click={scrollRight}>
            <IconArrowRight />
        </button>
    </div>
</div>
{/if}

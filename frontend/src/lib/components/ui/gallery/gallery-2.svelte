<script lang="ts">
    import { onMount, createEventDispatcher, tick } from "svelte";
    import type { ImageDetails, Loading } from './types';

    export let gap: number = 10;
    export let maxColumnWidth: number = 250;
    export let hover: boolean = false;
    export let loading: Loading = "eager";

    const dispatch = createEventDispatcher<{ click: ImageDetails }>();

    let slotHolder: HTMLElement | null = null;
    let columns: ImageDetails[][] = [];
    let galleryWidth: number = 0;
    let columnCount: number = 0;

    $: columnCount = Math.max(Math.floor(galleryWidth / maxColumnWidth), 1);
    $: galleryStyle = `grid-template-columns: repeat(${columnCount}, 1fr); --gap: ${gap}px`;

    $: {
        if (columnCount) {
            Draw();
        }
    }

    onMount(() => {
        Draw();

        const observer = new MutationObserver(Draw);
        if (slotHolder) {
            observer.observe(slotHolder, { childList: true, subtree: true });
        }

        return () => observer.disconnect();
    });

    function handleImageClick(img: ImageDetails) {
        dispatch("click", img);
    }

    async function Draw() {
        await tick();
        if (!slotHolder) return;

        const images = Array.from(slotHolder.querySelectorAll<HTMLImageElement>('img'));
        columns = Array.from({ length: columnCount }, () => []);

        images.forEach((img, i) => {
            columns[i % columnCount].push({ src: img.src, alt: img.alt, class: img.className, loading: loading as 'eager' | 'lazy' });
        });
    }
</script>

<div id="slotHolder" bind:this={slotHolder}>
    <slot />
</div>

{#if columns.length}
    <div id="gallery" bind:clientWidth={galleryWidth} style={galleryStyle}>
        {#each columns as column}
            <div class="column">
                {#each column as img}
                    <button
                        type="button"
                        on:click={() => handleImageClick(img)}
                        on:keydown={(e) => e.key === 'Enter' && handleImageClick(img)}
                        class="image-button {hover ? 'img-hover' : ''} {img.class}"
                    >
                    <img src={img.src} alt={img.alt} loading={img.loading} />
                </button>
                {/each}
            </div>
        {/each}
    </div>
{/if}

<style>
    #slotHolder {
        display: none;
    }
    #gallery {
        width: 100%;
        display: grid;
        gap: var(--gap);
    }
    #gallery .column {
        display: flex;
        flex-direction: column;
    }
    #gallery .column * {
        width: 100%;
        margin-top: var(--gap);
    }
    #gallery .column *:nth-child(1) {
        margin-top: 0;
    }
    .img-hover {
        opacity: 0.9;
        transition: all 0.2s;
    }
    .img-hover:hover {
        opacity: 1;
        transform: scale(1.05);
    }
</style>

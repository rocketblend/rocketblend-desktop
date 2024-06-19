<script lang="ts">
    import { createEventDispatcher, tick } from "svelte";
    import type { Loading, GalleryItem } from './types';
    import Media from './media.svelte';

    export let gap: number = 10;
    export let maxColumnWidth: number = 250;
    export let items: GalleryItem[] = [];
    export let highlight: string[] = [];
    export let loading: Loading = "eager";
    export let hover: boolean = false;
    export let rounded: boolean = false;

    const dispatch = createEventDispatcher<{ click: { value: string }, dblclick: { value: string, event: MouseEvent } }>();

    let columns: GalleryItem[][] = [];
    let galleryWidth: number = 0;
    let columnCount: number = 0;

    $: columnCount = Math.max(Math.floor(galleryWidth / maxColumnWidth), 1);
    $: galleryStyle = `grid-template-columns: repeat(${columnCount}, 1fr); --gap: ${gap}px`;

    $: {
        if (columnCount) {
            Draw();
        }
    }

    $: {
        if (items.length) {
            Draw();
        }
    }

    function handleClick(value: string) {
        dispatch("click", {value});
    }

    function handleDbClick(value: string, event: MouseEvent) {
        dispatch("dblclick", {value, event});
    }

    async function Draw() {
        await tick();

        columns = Array.from({ length: columnCount }, () => []);

        items.forEach((item, i) => {
            columns[i % columnCount].push(item);
        });
    }
</script>

{#if columns.length}
    <div id="gallery" bind:clientWidth={galleryWidth} style={galleryStyle}>
        {#each columns as column}
            <div class="column">
                {#each column as item, i}   
                    <div>
                        <Media
                            src={item.src}
                            alt={item.alt}
                            class={item.class}
                            highlight={highlight.includes(item.value)}
                            loading={loading}
                            hover={hover}
                            rounded={rounded}
                            on:click={(e) => handleClick(item.value)}
                            on:dblclick={(e) => handleDbClick(item.value, e)}
                        />
                    </div>
                {/each}
            </div>
        {/each}
    </div>
{/if}

<style>
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
</style>

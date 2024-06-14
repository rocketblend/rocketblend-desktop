<script lang="ts">
    import { onMount, createEventDispatcher, tick } from "svelte";
    import type { Loading, MediaDetails } from './types';
    import Media2 from './media-2.svelte';

    export let gap: number = 10;
    export let maxColumnWidth: number = 250;
    export let items: MediaDetails[] = [];
    export let loading: Loading = "eager";
    export let hover: boolean = false;
    export let rounded: boolean = false;

    const dispatch = createEventDispatcher<{ click: MediaDetails }>();

    let columns: MediaDetails[][] = [];
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
    });

    function handleImageClick(event: CustomEvent<MediaDetails>) {
        dispatch("click", event.detail);
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
                        <Media2
                            src={item.src}
                            alt={item.alt}
                            className={item.class}
                            loading={loading}
                            hover={hover}
                            rounded={rounded}
                            on:click={handleImageClick}
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
        padding-right: .25rem;
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

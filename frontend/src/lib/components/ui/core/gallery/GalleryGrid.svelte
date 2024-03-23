<script lang="ts">
    import type { MediaInfo } from '$lib/types';
    import Gallery from "./Gallery.svelte";

    export let items: MediaInfo[] = [];
    export let columnCount: number = 4;
    export let selectedIds: string[] = [];
    export let multiple: boolean = false;

    function distributeItems(items: MediaInfo[], columnCount: number): MediaInfo[][] {
        let columns: MediaInfo[][] = Array.from({ length: columnCount }, () => []);
        items.forEach((item, index) => {
            columns[index % columnCount].push(item);
        });

        return columns;
    }

    // Calculated galleries based on items and column count
    $: galleries = distributeItems(items, columnCount);
</script>

<Gallery class="gap-2 grid-cols-2 lg:grid-cols-4">
    {#each galleries as gallery}
        <Gallery multiple={multiple} items={gallery} bind:group={selectedIds} on:itemDoubleClick />
    {/each}
</Gallery>
<script lang="ts">
    import type { project } from '$lib/wailsjs/go/models';
    import type { MediaInfo } from '$lib/types';

    import { resourcePath } from '$lib/components/utils';
    
    import Gallery from "../core/gallery/Gallery.svelte";

    export let sourceData: project.Project[] = [];

    let columns: number = 4;
    let galleries: MediaInfo[][] = Array.from({ length: columns }, () => []);

    export let selectedIds: string[] = [];

    sourceData.forEach((proj, index) => {
        galleries[index % columns].push({
                id: proj.id?.toString() || "",
                title: proj.name || "",
                alt: `${proj.name || ""} splash`,
                src: resourcePath(proj.splashPath)
            }
        );
    });
</script>

<p>{selectedIds}</p>

<Gallery class="gap-2 grid-cols-2 lg:grid-cols-4">
    {#each galleries as galleryItems}
        <Gallery items={galleryItems} bind:group={selectedIds} multiple/>
    {/each}
</Gallery>
<script lang="ts">
    import type { types } from '$lib/wailsjs/go/models';

    import { DisplayType } from '$lib/types';

    import { Gallery, type GalleryItem } from '$lib/components/ui/gallery';
    
    import ProjectTable from "./project-table.svelte";

    export let projects: types.Project[] = [];
    export let selectedProjectIds: string[] = [];
    export let displayType: DisplayType = DisplayType.Table;

    let galleryItems: GalleryItem[] = [];

    function handleGalleryClick(event: CustomEvent<{value: string}>) {
        const { value } = event.detail;
        selectedProjectIds = [value];
    }

    function convertProjectsToGalleryItems(projects: types.Project[] = []): GalleryItem[] {
        return projects.map((project) => ({
            value: project.id.toString(),
            src: project.splash?.url || "",
            alt: `${project.name || ""}`,
            class: "",
        }));
    }

    $: galleryItems = convertProjectsToGalleryItems(projects);
</script>

<div class="overflow-auto h-full p-1">
    {#if projects === undefined || projects.length === 0}
        <div class="flex items-center justify-center h-64">
            <h4>No projects found.</h4>
        </div>
    {:else}
        {#if displayType === DisplayType.Gallery}
            <Gallery
                gap={15}
                maxColumnWidth={250}
                hover={true}
                bind:items={galleryItems}
                bind:highlight={selectedProjectIds}
                loading="eager"
                rounded={true}
                on:click={handleGalleryClick}
                on:dblclick
            />
        {:else}
            <ProjectTable
                on:dblclick
                bind:sourceData={projects}
                bind:selectedProjectIds={selectedProjectIds} />
        {/if}
    {/if}
</div>

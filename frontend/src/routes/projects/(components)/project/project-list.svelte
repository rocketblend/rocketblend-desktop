<script lang="ts">
    import type { types } from '$lib/wailsjs/go/models';

    import { DisplayType, type MediaInfo } from '$lib/types';
    import { GalleryGrid } from '$lib/components/ui/gallery';
    
    import ProjectTable from "./project-table.svelte";

    export let projects: types.Project[] = [];
    export let selectedProjectIds: string[] = [];
    export let displayType: DisplayType = DisplayType.Table;

    function convertProjectsToGalleryItems(projects: types.Project[] = []): MediaInfo[] {
        return projects.map((project) => ({
            id: project.id?.toString() || "",
            title: project.name || "",
            alt: `${project.name || ""} splash`,
            src: project.splash?.url || "",
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
            <GalleryGrid
                on:itemDoubleClick
                bind:items={galleryItems}
                bind:selectedIds={selectedProjectIds}/>
        {:else }
            <ProjectTable
                on:sortChanged
                on:itemDoubleClick
                bind:sourceData={projects}
                bind:selectedProjectIds={selectedProjectIds} />
        {/if}
    {/if}
</div>

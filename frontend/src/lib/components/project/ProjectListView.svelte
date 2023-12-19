<script lang="ts">
    import ProjectTable from '$lib/components/project/ProjectTable.svelte';
    import GalleryGrid from '$lib/components/core/gallery/GalleryGrid.svelte';
    import type { project } from '$lib/wailsjs/go/models';
    import { DisplayType, type MediaInfo } from '$lib/types';
    import { resourcePath } from '$lib/components/utils';

    export let projects: project.Project[] = [];
    export let selectedProjectIds: string[] = [];
    export let displayType: DisplayType = DisplayType.Table;

    function convertProjectsToGalleryItems(projects: project.Project[] = []): MediaInfo[] {
        return projects.map((project) => ({
            id: project.id?.toString() || "",
            title: project.name || "Untitled Project",
            alt: `${project.name || "Untitled Project"} splash`,
            src: resourcePath(project.splashPath)
        }));
    }

    $: galleryItems = convertProjectsToGalleryItems(projects);
</script>

{#if projects === undefined || projects.length === 0}
    <div class="flex items-center justify-center h-64">
        <h4>No projects found.</h4>
    </div>
{:else}
    {#if displayType === DisplayType.Gallery}
        <GalleryGrid
            on:itemDoubleClicked
            on:ctrlItemDoubleClicked
            bind:items={galleryItems}
            bind:selectedIds={selectedProjectIds}/>
    {:else }
        <ProjectTable bind:sourceData={projects} on:selected />
    {/if}
{/if}
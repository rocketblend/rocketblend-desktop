<script lang="ts">
    import type { types } from '$lib/wailsjs/go/models';

    import { DisplayType, type MediaInfo } from '$lib/types';
    import { Gallery2, GalleryGrid } from '$lib/components/ui/gallery';
    
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
            <Gallery2
                gap={15}
                maxColumnWidth={250}
                hover={true}
                loading="lazy"
            >
                <img
                    src="https://via.placeholder.com/210x170/100"
                    alt="210x170"
                    class="hi"
                />
                <img src="https://via.placeholder.com/180x200/100" alt="180x200" />
                <img src="https://via.placeholder.com/200x210/100" alt="200x210" />
                <img src="https://via.placeholder.com/140x250/100" alt="140x250" />
                <img src="https://via.placeholder.com/250x300/100" alt="250x300" />
                <img src="https://via.placeholder.com/280x200/100" alt="280x200" />
                <img src="https://via.placeholder.com/220x180/100" alt="220x180" />
                <img src="https://via.placeholder.com/180x150/100" alt="180x150" />
                <img src="https://via.placeholder.com/210x210/100" alt="210x210" />
                <img src="https://via.placeholder.com/200x200/100" alt="200x200" />
                <img src="https://via.placeholder.com/220x200/100" alt="220x200" />
                <img src="https://via.placeholder.com/180x310/100" alt="180x310" />
                <img src="https://via.placeholder.com/210x210/100" alt="210x210" />
                <img src="https://via.placeholder.com/200x280/100" alt="200x280" />
                <img src="https://via.placeholder.com/210x350/100" alt="210x350" />
                <img src="https://via.placeholder.com/180x270/100" alt="180x270" />
            </Gallery2>
            <!-- <GalleryGrid
                on:itemDoubleClick
                bind:items={galleryItems}
                bind:selectedIds={selectedProjectIds}/> -->
        {:else }
            <ProjectTable
                on:sortChanged
                on:itemDoubleClick
                bind:sourceData={projects}
                bind:selectedProjectIds={selectedProjectIds} />
        {/if}
    {/if}
</div>

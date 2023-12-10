<script lang="ts">
    import { goto } from '$app/navigation';
    import { page } from '$app/stores';
    import { RadioGroup, RadioItem } from '@skeletonlabs/skeleton';
    import { t } from '$lib/translations/translations';
    import SearchInput from '$lib/components/core/input/SearchInput.svelte';
    import ProjectTable from '$lib/components/project/ProjectTable.svelte';
    import GalleryGrid from '$lib/components/core/gallery/GalleryGrid.svelte';
    import { resourcePath } from '$lib/components/utils';
    import { selectedProjectIds } from '$lib/store';
    import { RunProject } from '$lib/wailsjs/go/application/Driver';
    import type { project } from '$lib/wailsjs/go/models';
    import type { MediaInfo } from '$lib/types';
    import type { PageData } from './$types';

    export let data: PageData;

    let searchQuery = '';
    let displayType = 'table';
    let form : HTMLFormElement;

    function handleSelected(event: CustomEvent<project.Project | null>): void {
        const project = event.detail;
        if (!project || !project.id) {
            return;
        }

        selectedProjectIds.set([project.id.toString()]);
    }

    function handleProjectDoubleClick(event: CustomEvent<{ itemId: string }>) {
        goto(`/projects/${event.detail.itemId}`);
    }

    function handleProjectActionDoubleClick(event: CustomEvent<{ itemId: string }>) {
        RunProject(event.detail.itemId)
    }

    function handleFormSubmit(event: Event): void {
        form.requestSubmit();
    }

    function convertProjectsToGalleryItems(projects: project.Project[]): MediaInfo[] {
        return projects.map((project) => ({
            id: project.id?.toString() || "",
            title: project.name || "Untitled Project",
            alt: `${project.name || "Untitled Project"} splash`,
            src: resourcePath(project.splashPath)
        }));
    }

    $: searchQuery = $page.url.searchParams.get('query') || '';
    $: displayType = $page.url.searchParams.get('display') || 'table';
    $: galleryItems = convertProjectsToGalleryItems(data.projects || []);
</script>

<main class="space-y-4"> 
    <h2 class="font-bold">{$t('home.title')}</h2>
    <div class="space-y-4">
        <form bind:this={form} data-sveltekit-keepfocus class="inline-flex space-x-4 w-full">
            <div class="flex-grow">
                <SearchInput name="query" value={searchQuery} on:input={handleFormSubmit} placeholder="Search" debounceDelay={500} />
            </div>

            <RadioGroup>
                <RadioItem name="display" group={displayType} value={"table"} on:change={handleFormSubmit}>Table</RadioItem>
                <RadioItem name="display" group={displayType} value={"gallery"} on:change={handleFormSubmit}>Gallery</RadioItem>
            </RadioGroup>
            <button type="submit" class="hidden">Search</button>
        </form>
        
        {#if data.projects === undefined || data.projects.length === 0}
            <div class="flex items-center justify-center h-64">
                <h4>No projects found.</h4>
            </div>
        {:else}
            {#if displayType === "gallery"}
                <GalleryGrid
                    on:itemDoubleClicked={handleProjectDoubleClick}
                    on:ctrlItemDoubleClicked={handleProjectActionDoubleClick}
                    bind:items={galleryItems}
                    bind:selectedIds={$selectedProjectIds}/>
            {:else }
                <ProjectTable bind:sourceData={data.projects} on:selected={handleSelected} />
            {/if}
        {/if}
    </div>
</main>

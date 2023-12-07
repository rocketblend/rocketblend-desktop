<script lang="ts">
    import { t } from '$lib/translations/translations';

    import type { PageData } from './$types';
    import { page } from '$app/stores'

    import { RadioGroup, RadioItem } from '@skeletonlabs/skeleton';
    
    import SearchInput from '$lib/components/core/input/SearchInput.svelte';
	import ProjectGallery from '$lib/components/project/ProjectGallery.svelte';
    import ProjectTable from '$lib/components/project/ProjectTable.svelte';

    import type { project } from '$lib/wailsjs/go/models';

    import { selectedProjectIds } from '$lib/store';

    export let data: PageData;

    let searchQuery = '';
    let displayType = 'table';
    let form : HTMLFormElement;

    $: searchQuery = $page.url.searchParams.get('query') || '';
    $: displayType = $page.url.searchParams.get('display') || 'table';

    function handleSelected(event: CustomEvent<project.Project | null>): void {
        const project = event.detail;
        if (!project || !project.id) {
            return;
        }

        selectedProjectIds.set([project.id.toString()]);
    }

    function handleFormSubmit(event: Event): void {
        form.requestSubmit();
    }
</script>

<main class="space-y-4"> 
    <h2 class="font-bold">{$t('home.title')}</h2>
    <div class="space-y-4">
        <div class="flex items-center justify-between space-x-4">
            <div class="w-full">
                <form bind:this={form} data-sveltekit-keepfocus>
                    <SearchInput name="query" value={searchQuery} on:input={handleFormSubmit} placeholder="Search" debounceDelay={500}/>
                    <RadioGroup>
                        <RadioItem name="display" group={displayType} value={"table"} on:change={handleFormSubmit}>Table</RadioItem>
                        <RadioItem name="display" group={displayType} value={"gallery"} on:change={handleFormSubmit}>Gallery</RadioItem>
                    </RadioGroup>
                    <button type="submit" class="hidden">Search</button>
                </form>
            </div>
        </div>
        
        {#if data.projects === undefined || data.projects.length === 0}
            <div class="flex items-center justify-center h-64">
                <h4>No projects found.</h4>
            </div>
        {:else}
            {#if displayType === "gallery"}
                <ProjectGallery bind:sourceData={data.projects} bind:selectedIds={$selectedProjectIds}/>
            {:else }
                <ProjectTable bind:sourceData={data.projects} on:selected={handleSelected} />
            {/if}
        {/if}
    </div>
</main>

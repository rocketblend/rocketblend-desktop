<script lang="ts">
    import { onMount } from 'svelte';

    import { RadioGroup, RadioItem } from '@skeletonlabs/skeleton';
    
	import ProjectGallery from '../project-gallery/project-gallery.svelte';
    import ProjectTable from '../project-table/project-table.svelte';
    import SearchInput from '../core/search-input/search-input.svelte';

    import type { projectservice } from '$lib/wailsjs/go/models';
    import { FindAllProjects } from '$lib/wailsjs/go/application/Driver';

    let projects: projectservice.Project[] = [];
    let searchTerm: string = '';
    let displayType: number = 0;

    async function loadProjects(query: string): Promise<void> {
        const result = await FindAllProjects(query);
        projects = result || [];
    }

    function handleSearch(event: CustomEvent<string>): void {
        searchTerm = event.detail;
    }

    onMount(async () => {
        await loadProjects(searchTerm);
    });

    $: loadProjects(searchTerm), searchTerm;
</script>

<div class="space-y-4">
    <div class="flex items-center justify-between space-x-4">
        <div class="w-full">
            <SearchInput on:search={handleSearch} placeholder="Search"/>
        </div>
        <RadioGroup>
            <RadioItem bind:group={displayType} name="justify" value={0}>Table</RadioItem>
            <RadioItem bind:group={displayType} name="justify" value={1}>Gallery</RadioItem>
        </RadioGroup>
    </div>
    
    {#if projects.length === 0}
        <div class="flex items-center justify-center h-64">
            <h4>No projects found.</h4>
        </div>
    {:else}
        {#if displayType === 0}
            <ProjectTable sourceData={projects}/>
        {:else if displayType === 1}
            <ProjectGallery />
        {/if}
    {/if}

</div>
